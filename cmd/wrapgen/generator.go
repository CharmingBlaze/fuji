package main

import (
	"fmt"
	"go/token"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// API represents the extracted API from C/C++ headers
type API struct {
	Name         string
	Functions    []Function
	Structs      []Struct
	Enums        []Enum
	Macros       []Macro
	Typedefs     []Typedef
	Constants    []Constant
	Headers      []string
	Dependencies []string
}

type Function struct {
	Name          string
	ReturnType    string
	Parameters    []Parameter
	Variadic      bool
	Static        bool
	Inline        bool
	Documentation string
	Header        string
}

type Parameter struct {
	Name    string
	Type    string
	Default string
}

type Struct struct {
	Name          string
	Fields        []Field
	Methods       []Function
	Documentation string
	Header        string
}

type Field struct {
	Name string
	Type string
}

type Enum struct {
	Name          string
	Values        []EnumValue
	Documentation string
	Header        string
}

type EnumValue struct {
	Name  string
	Value int64
}

type Macro struct {
	Name          string
	Value         string
	Parameters    []string
	Documentation string
	Header        string
}

type Typedef struct {
	Name          string
	TargetType    string
	Documentation string
	Header        string
}

type Constant struct {
	Name          string
	Type          string
	Value         string
	Documentation string
	Header        string
}

// wrapgenBlockedBindingNames are not all Go keywords but cannot be used as Fuji `let` names
// (regex false-positives from C headers / macros, or C reserved words).
var wrapgenBlockedBindingNames = map[string]struct{}{
	"let": {}, "func": {}, "import": {}, "export": {},
	"if": {}, "else": {}, "while": {}, "for": {}, "switch": {}, "case": {}, "default": {},
	"break": {}, "continue": {}, "return": {}, "goto": {}, "do": {},
	"sizeof": {}, "typeof": {}, "asm": {}, "typeof_unqual": {},
	"struct": {}, "union": {}, "enum": {}, "typedef": {},
	"signed": {}, "unsigned": {}, "const": {}, "volatile": {},
	"static": {}, "extern": {}, "inline": {}, "restrict": {},
	"auto": {}, "register": {}, "alignas": {}, "alignof": {},
	"bool": {}, "true": {}, "false": {},
	"void": {}, "int": {}, "char": {}, "short": {}, "long": {}, "float": {}, "double": {},
}

// isValidWrapgenBindingName rejects regex false-positives (e.g. macro fragments parsed as `if(...)`) and names that cannot be emitted as Fuji identifiers.
func isValidWrapgenBindingName(name string) bool {
	if name == "" {
		return false
	}
	r, w := utf8.DecodeRuneInString(name)
	if w == 0 || r == utf8.RuneError {
		return false
	}
	if r != '_' && !unicode.IsLetter(r) {
		return false
	}
	for _, ch := range name[w:] {
		if ch != '_' && !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}
	if token.IsKeyword(name) {
		return false
	}
	if _, bad := wrapgenBlockedBindingNames[name]; bad {
		return false
	}
	return true
}

func filterAndDedupeFunctions(fs []Function) []Function {
	seen := make(map[string]struct{})
	var out []Function
	for _, f := range fs {
		if !isValidWrapgenBindingName(f.Name) {
			continue
		}
		if _, ok := seen[f.Name]; ok {
			continue
		}
		seen[f.Name] = struct{}{}
		out = append(out, f)
	}
	return out
}

type WrapperGenerator struct {
	config *WrapGenConfig
	fset   *token.FileSet
}

func NewWrapperGenerator(config *WrapGenConfig) *WrapperGenerator {
	return &WrapperGenerator{
		config: config,
		fset:   token.NewFileSet(),
	}
}

// ParseHeaders extracts API information from C/C++ header files
func (wg *WrapperGenerator) ParseHeaders() (*API, error) {
	api := &API{
		Name: wg.config.LibraryName,
	}

	// Use enhanced regex parsing for better compatibility with complex headers
	if wg.config.Verbose {
		fmt.Printf("Using enhanced regex parser for comprehensive C/C++ analysis...\n")
	}

	// Parse all headers with enhanced regex
	for _, header := range wg.config.InputHeaders {
		if wg.config.Verbose {
			fmt.Printf("Parsing header: %s\n", header)
		}

		headerAPI, err := wg.parseHeader(header)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %v", header, err)
		}

		// Merge header API into main API
		api.Functions = append(api.Functions, headerAPI.Functions...)
		api.Structs = append(api.Structs, headerAPI.Structs...)
		api.Enums = append(api.Enums, headerAPI.Enums...)
		api.Macros = append(api.Macros, headerAPI.Macros...)
		api.Typedefs = append(api.Typedefs, headerAPI.Typedefs...)
		api.Constants = append(api.Constants, headerAPI.Constants...)
		api.Headers = append(api.Headers, header)
	}

	api.Functions = filterAndDedupeFunctions(api.Functions)

	// Analyze dependencies and relationships
	wg.analyzeDependencies(api)

	if wg.config.Verbose {
		fmt.Printf("Found %d functions, %d structs, %d enums, %d macros, %d typedefs, %d constants\n",
			len(api.Functions), len(api.Structs), len(api.Enums), len(api.Macros), len(api.Typedefs), len(api.Constants))
	}

	return api, nil
}

// parseHeader parses a single C/C++ header file
func (wg *WrapperGenerator) parseHeader(headerPath string) (*API, error) {
	content, err := ioutil.ReadFile(headerPath)
	if err != nil {
		return nil, err
	}

	headerContent := string(content)
	// Strip all comments for reliable regex matching
	reMulti := regexp.MustCompile(`(?s)/\*.*?\*/`)
	headerContent = reMulti.ReplaceAllString(headerContent, "")
	reSingle := regexp.MustCompile(`//.*`)
	headerContent = reSingle.ReplaceAllString(headerContent, "")

	api := &API{}

	// Extract functions
	api.Functions = wg.extractFunctions(headerContent, headerPath)

	// Extract structs
	api.Structs = wg.extractStructs(headerContent, headerPath)

	// Extract enums
	api.Enums = wg.extractEnums(headerContent, headerPath)

	// Extract macros
	api.Macros = wg.extractMacros(headerContent, headerPath)

	// Extract typedefs
	api.Typedefs = wg.extractTypedefs(headerContent, headerPath)

	// Extract constants
	api.Constants = wg.extractConstants(headerContent, headerPath)

	return api, nil
}

// extractFunctions extracts function declarations from header content
func (wg *WrapperGenerator) extractFunctions(content, headerPath string) []Function {
	var functions []Function

	// Function regex pattern
	funcPattern := regexp.MustCompile(`(?m)^\s*(?:inline\s+)?(?:static\s+)?(?:\w+\s+)*?(\w+)\s*\(([^)]*)\)\s*(?:__attribute__\s*\([^)]*\))?\s*;?`)

	matches := funcPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			funcName := strings.TrimSpace(match[1])
			paramsStr := strings.TrimSpace(match[2])

			if !isValidWrapgenBindingName(funcName) {
				continue
			}

			// Skip if this looks like a struct/union definition
			if wg.isStructDefinition(funcName, paramsStr) {
				continue
			}

			function := Function{
				Name:   funcName,
				Header: headerPath,
			}

			// Extract return type (simplified)
			fullMatch := match[0]
			returnType := wg.extractReturnType(fullMatch, funcName)
			function.ReturnType = returnType

			// Parse parameters
			function.Parameters = wg.parseParameters(paramsStr)

			functions = append(functions, function)
		}
	}

	return functions
}

// extractStructs extracts struct definitions from header content
func (wg *WrapperGenerator) extractStructs(content, headerPath string) []Struct {
	var structs []Struct

	// Pattern 1: typedef struct [Name] { ... } [TypedefName];
	structPattern1 := regexp.MustCompile(`(?s)typedef\s+struct\s*(\w*)\s*\{([^}]+)\}\s*(\w+)[^;]*;`)
	// Pattern 2: struct Name { ... };
	structPattern2 := regexp.MustCompile(`(?s)struct\s+(\w+)\s*\{([^}]+)\}\s*;`)

	patterns := []struct {
		re   *regexp.Regexp
		nIdx int
		fIdx int
	}{
		{structPattern1, 3, 2},
		{structPattern2, 1, 2},
	}

	seen := make(map[string]bool)

	for _, p := range patterns {
		matches := p.re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > p.nIdx && len(match) > p.fIdx {
				structName := strings.TrimSpace(match[p.nIdx])
				if structName == "" && p.nIdx == 3 {
					structName = strings.TrimSpace(match[1]) // fallback to tag name
				}

				if structName == "" || seen[structName] {
					continue
				}
				seen[structName] = true

				structDef := Struct{
					Name:   structName,
					Header: headerPath,
				}

				// Parse fields
				fieldsContent := match[p.fIdx]
				structDef.Fields = wg.parseStructFields(fieldsContent)

				structs = append(structs, structDef)
			}
		}
	}

	return structs
}

// extractEnums extracts enum definitions from header content
func (wg *WrapperGenerator) extractEnums(content, headerPath string) []Enum {
	var enums []Enum

	// Enum regex patterns
	enumPattern1 := regexp.MustCompile(`(?s)typedef\s+enum\s+(\w*)\s*\{([^}]+)\}\s*(\w+)[^;]*;`)
	enumPattern2 := regexp.MustCompile(`(?s)enum\s+(\w+)\s*\{([^}]+)\}\s*;`)

	patterns := []*regexp.Regexp{enumPattern1, enumPattern2}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				enumName := strings.TrimSpace(match[3])
				if enumName == "" {
					enumName = strings.TrimSpace(match[1])
				}

				enumDef := Enum{
					Name:   enumName,
					Header: headerPath,
				}

				// Parse enum values
				valuesContent := match[2]
				enumDef.Values = wg.parseEnumValues(valuesContent)

				enums = append(enums, enumDef)
			}
		}
	}

	return enums
}

// extractMacros extracts macro definitions from header content
func (wg *WrapperGenerator) extractMacros(content, headerPath string) []Macro {
	var macros []Macro

	// Macro regex pattern
	macroPattern := regexp.MustCompile(`(?m)^#define\s+(\w+)\s*(?:\(([^)]*)\))?\s*(.+)$`)

	matches := macroPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			macroName := strings.TrimSpace(match[1])
			paramsStr := ""
			value := ""

			if len(match) >= 4 && match[2] != "" {
				// Function-like macro
				paramsStr = strings.TrimSpace(match[2])
				value = strings.TrimSpace(match[3])
			} else if len(match) >= 3 {
				// Object-like macro
				value = strings.TrimSpace(match[2])
			}

			macro := Macro{
				Name:   macroName,
				Value:  value,
				Header: headerPath,
			}

			if paramsStr != "" {
				macro.Parameters = strings.Split(paramsStr, ",")
				for i, param := range macro.Parameters {
					macro.Parameters[i] = strings.TrimSpace(param)
				}
			}

			macros = append(macros, macro)
		}
	}

	return macros
}

// extractTypedefs extracts typedef declarations from header content
func (wg *WrapperGenerator) extractTypedefs(content, headerPath string) []Typedef {
	var typedefs []Typedef

	// Typedef regex pattern (excluding function pointers and complex blocks for now)
	// Supports: typedef TargetType NewName;
	typedefPattern := regexp.MustCompile(`typedef\s+([a-zA-Z_][a-zA-Z0-9_\s\*]+)\s+([a-zA-Z_][a-zA-Z0-9_]+);`)

	matches := typedefPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			targetType := strings.TrimSpace(match[1])
			typeName := strings.TrimSpace(match[2])

			typedef := Typedef{
				Name:       typeName,
				TargetType: targetType,
				Header:     headerPath,
			}

			typedefs = append(typedefs, typedef)
		}
	}

	return typedefs
}

// extractConstants extracts constant definitions from header content
func (wg *WrapperGenerator) extractConstants(content, headerPath string) []Constant {
	var constants []Constant

	// Constant regex patterns
	constPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^extern\s+const\s+(\w+)\s+(\w+)[^;]*;`),
		regexp.MustCompile(`(?m)^const\s+(\w+)\s+(\w+)\s*=\s*([^;]+);`),
	}

	for _, pattern := range constPatterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				constType := strings.TrimSpace(match[1])
				constName := strings.TrimSpace(match[2])
				constValue := ""

				if len(match) >= 4 {
					constValue = strings.TrimSpace(match[3])
				}

				constant := Constant{
					Name:   constName,
					Type:   constType,
					Value:  constValue,
					Header: headerPath,
				}

				constants = append(constants, constant)
			}
		}
	}

	return constants
}

// Helper methods for parsing specific components
func (wg *WrapperGenerator) extractReturnType(fullMatch, funcName string) string {
	// Extract return type by removing function name and parameters
	beforeFunc := strings.Split(fullMatch, funcName)[0]
	returnType := strings.TrimSpace(beforeFunc)

	// Clean up return type
	returnType = strings.ReplaceAll(returnType, "inline", "")
	returnType = strings.ReplaceAll(returnType, "static", "")
	returnType = strings.ReplaceAll(returnType, "RLAPI", "")
	returnType = strings.TrimSpace(returnType)

	return returnType
}

func (wg *WrapperGenerator) parseParameters(paramsStr string) []Parameter {
	var parameters []Parameter

	if paramsStr == "" || paramsStr == "void" {
		return parameters
	}

	paramPairs := strings.Split(paramsStr, ",")
	for _, param := range paramPairs {
		param = strings.TrimSpace(param)
		if param == "" || param == "void" {
			continue
		}

		parts := strings.Fields(param)
		if len(parts) >= 2 {
			paramName := parts[len(parts)-1]
			// Handle pointer attached to name: char *name, char **name
			for strings.HasPrefix(paramName, "*") {
				paramName = paramName[1:]
				if len(parts) >= 2 {
					parts[len(parts)-2] += "*"
				}
			}
			paramType := strings.Join(parts[:len(parts)-1], " ")

			// Handle default values
			if strings.Contains(paramName, "=") {
				nameValue := strings.Split(paramName, "=")
				paramName = strings.TrimSpace(nameValue[0])
				// defaultValue = strings.TrimSpace(nameValue[1])
			}

			parameters = append(parameters, Parameter{
				Name: paramName,
				Type: paramType,
			})
		}
	}

	return parameters
}

func (wg *WrapperGenerator) parseStructFields(fieldsContent string) []Field {
	var fields []Field

	// Strip comments from the entire block first
	reMulti := regexp.MustCompile(`(?s)/\*.*?\*/`)
	fieldsContent = reMulti.ReplaceAllString(fieldsContent, "")
	reSingle := regexp.MustCompile(`//.*`)
	fieldsContent = reSingle.ReplaceAllString(fieldsContent, "")

	lines := strings.Split(fieldsContent, ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			fieldName := parts[len(parts)-1]
			// Handle pointer fields: char *name, char **name
			for strings.HasPrefix(fieldName, "*") {
				fieldName = fieldName[1:]
				if len(parts) >= 2 {
					parts[len(parts)-2] += "*"
				}
			}
			fieldType := strings.Join(parts[:len(parts)-1], " ")

			fields = append(fields, Field{
				Name: fieldName,
				Type: fieldType,
			})
		}
	}

	return fields
}

func (wg *WrapperGenerator) isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}

func (wg *WrapperGenerator) parseEnumValues(valuesContent string) []EnumValue {
	var values []EnumValue

	// Replace newlines with spaces to handle multi-line enums
	valuesContent = strings.ReplaceAll(valuesContent, "\n", " ")
	valuesContent = strings.ReplaceAll(valuesContent, "\r", " ")

	lines := strings.Split(valuesContent, ",")
	currentValue := int64(0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove comments
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		if idx := strings.Index(line, "/*"); idx != -1 {
			if endIdx := strings.Index(line, "*/"); endIdx != -1 {
				line = line[:idx] + line[endIdx+2:]
			} else {
				line = line[:idx]
			}
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")
		valueName := strings.TrimSpace(parts[0])
		if valueName == "" {
			continue
		}

		// Final cleanup of valueName (remove any remaining non-identifier chars)
		partsName := strings.Fields(valueName)
		if len(partsName) > 0 {
			valueName = partsName[len(partsName)-1]
		}

		if !wg.isValidIdentifier(valueName) {
			continue
		}

		if len(parts) > 1 {
			valueStr := strings.TrimSpace(parts[1])
			if strings.HasPrefix(valueStr, "0x") {
				fmt.Sscanf(valueStr, "%x", &currentValue)
			} else {
				fmt.Sscanf(valueStr, "%d", &currentValue)
			}
		}

		values = append(values, EnumValue{
			Name:  valueName,
			Value: currentValue,
		})
		currentValue++
	}

	return values
}

func (wg *WrapperGenerator) isStructDefinition(name, params string) bool {
	// Heuristic to detect if this is actually a struct definition
	return strings.Contains(params, "{") || strings.Contains(name, "struct")
}

// GenerateWrapper generates the wrapper code for the target language
func (wg *WrapperGenerator) GenerateWrapper(api *API) error {
	switch wg.config.Language {
	case "fuji", "":
		return wg.generateFujiWrapper(api)
	default:
		return fmt.Errorf("unsupported language %q - only fuji (.fuji) is supported", wg.config.Language)
	}
}

// generateFujiWrapper generates an elegant, readable Fuji wrapper library.
func (wg *WrapperGenerator) generateFujiWrapper(api *API) error {
	wrapperFile := filepath.Join(wg.config.OutputDir, api.Name+".fuji")

	if wg.config.Verbose {
		fmt.Printf("  writing %s\n", wrapperFile)
	}

	file, err := os.Create(wrapperFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// ------------------------------------------------------------------ header
	w := func(format string, args ...interface{}) {
		fmt.Fprintf(file, format, args...)
	}

	w("// ----------------------------------------------------------------------------\n")
	w("// %s\n", api.Name)
	w("// Fuji library bindings - generated by wrapgen\n")
	w("// Ship this file with your app bundle.\n")
	w("//\n")
	w("// Usage in your Fuji program:\n")
	w("//   #include \"%s.fuji\"\n", api.Name)
	w("//\n")
	w("// Build and bundle:\n")
	w("//   set FUJI_NATIVE_SOURCES=wrapper.c\n")
	w("//   set FUJI_LINKFLAGS=-I<inc> -L<lib> -l%s\n", api.Name)
	w("//   fuji build mygame.fuji -o mygame.exe\n")
	w("//   fuji bundle mygame.fuji -o dist/mygame\n")
	w("// ----------------------------------------------------------------------------\n\n")

	// ---------------------------------------------------------------- constants
	if len(api.Constants) > 0 {
		w("// ----------------------------------------------------------------------------\n")
		w("// Constants\n")
		w("// ----------------------------------------------------------------------------\n\n")
		for _, constant := range api.Constants {
			wg.writeFujiConstant(file, constant)
		}
		w("\n")
	}

	// ------------------------------------------------------------------ structs
	if len(api.Structs) > 0 {
		w("// ----------------------------------------------------------------------------\n")
		w("// Structs\n")
		w("// ----------------------------------------------------------------------------\n\n")
		for _, structDef := range api.Structs {
			wg.writeFujiStructDefinition(file, structDef)
		}
		w("\n")
	}

	// ------------------------------------------------------------------- enums
	if len(api.Enums) > 0 {
		w("// ----------------------------------------------------------------------------\n")
		w("// Enums\n")
		w("// ----------------------------------------------------------------------------\n\n")
		for _, enum := range api.Enums {
			wg.writeFujiEnumDefinition(file, enum)
		}
		w("\n")
	}

	// ---------------------------------------------------------------- functions
	if len(api.Functions) > 0 {
		w("// ----------------------------------------------------------------------------\n")
		w("// Functions\n")
		w("// ----------------------------------------------------------------------------\n")
		for _, function := range api.Functions {
			wg.writeFujiFunctionDeclaration(file, function)
		}
	}

	return wg.generateCGlue(api)
}

func (wg *WrapperGenerator) writeFujiFunctionDeclaration(file *os.File, function Function) {
	ret := strings.TrimSpace(function.ReturnType)
	params := wg.bindingParamList(function)

	fmt.Fprintf(file, "\n")

	// Doc comment with signature and param list
	if ret == "" || ret == "void" {
		fmt.Fprintf(file, "// %s(%s)\n", function.Name, params)
	} else {
		fmt.Fprintf(file, "// %s(%s) -> %s\n", function.Name, params, ret)
	}
	for i, param := range function.Parameters {
		name := strings.TrimSpace(param.Name)
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		fmt.Fprintf(file, "//   @param %s\n", name)
	}

	// The extern directive (compiler directive, kept on its own line)
	fmt.Fprintf(file, "// fuji:extern %s %s %d\n", function.Name, wg.wrapperSymbol(function), len(function.Parameters))
	fmt.Fprintf(file, "let %s = 0;\n", function.Name)
}

func (wg *WrapperGenerator) bindingParamList(function Function) string {
	names := make([]string, 0, len(function.Parameters))
	for i, param := range function.Parameters {
		name := strings.TrimSpace(param.Name)
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

func (wg *WrapperGenerator) wrapperSymbol(function Function) string {
	return "fuji_wrap_" + sanitizeIdent(wg.config.LibraryName) + "_" + sanitizeIdent(function.Name)
}

func sanitizeIdent(s string) string {
	var b strings.Builder
	for i, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (i > 0 && r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return "wrapper"
	}
	return b.String()
}

func (wg *WrapperGenerator) generateCGlue(api *API) error {
	path := filepath.Join(wg.config.OutputDir, "wrapper.c")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "/*\n")
	fmt.Fprintf(file, " * %s — C library glue for the Fuji/Fuji native ABI\n", api.Name)
	fmt.Fprintf(file, " * Auto-generated by wrapgen. Do not edit.\n")
	fmt.Fprintf(file, " *\n")
	fmt.Fprintf(file, " * Compiled with your app via FUJI_NATIVE_SOURCES (e.g. fuji build / fuji bundle).\n")
	fmt.Fprintf(file, " * Each export wraps the C library: Value fn(int argCount, Value* args) -> native calls.\n")
	fmt.Fprintf(file, " */\n\n")
	fmt.Fprintf(file, "#include \"fuji_wrapgen_abi.h\"\n")
	fmt.Fprintf(file, "#include <stdbool.h>\n")
	fmt.Fprintf(file, "#include <string.h>\n")
	for _, h := range api.Headers {
		fmt.Fprintf(file, "#include \"%s\"\n", filepath.Base(h))
	}
	fmt.Fprintf(file, "\n")
	seen := make(map[string]bool)
	for _, function := range api.Functions {
		if function.Variadic {
			continue
		}
		if wg.skipGlueFunction(function) {
			continue
		}
		sym := wg.wrapperSymbol(function)
		if seen[sym] {
			continue
		}
		seen[sym] = true
		wg.writeCGlueFunction(file, api, function)
	}
	return nil
}

func (wg *WrapperGenerator) findStruct(api *API, typeName string) *Struct {
	typeName = strings.TrimSpace(typeName)
	// Remove 'struct ' prefix if present
	typeName = strings.TrimPrefix(typeName, "struct ")

	// Direct search
	for _, s := range api.Structs {
		if s.Name == typeName {
			return &s
		}
	}

	// Try finding via typedef
	for _, td := range api.Typedefs {
		if td.Name == typeName {
			// Resolve the underlying type
			underlying := strings.TrimSpace(td.TargetType)
			underlying = strings.TrimPrefix(underlying, "struct ")
			underlying = strings.TrimPrefix(underlying, "typedef ")
			for _, s := range api.Structs {
				if s.Name == underlying {
					return &s
				}
			}
		}
	}

	return nil
}

func (wg *WrapperGenerator) skipGlueFunction(function Function) bool {
	n := strings.TrimSpace(function.Name)
	if n == "" || strings.ContainsAny(n, "() ") {
		return true
	}
	// Regex parser sometimes emits bogus decls (e.g. typedefs mistaken for functions).
	if n == "void" || n == "bool" || n == "const" {
		return true
	}
	return false
}

func (wg *WrapperGenerator) writeCGlueFunction(file *os.File, api *API, function Function) {
	retType := strings.TrimSpace(function.ReturnType)
	retType = strings.TrimPrefix(retType, "RLAPI ")
	retType = strings.TrimSpace(retType)

	fmt.Fprintf(file, "/* %s(%s) -> %s */\n", function.Name, wg.cParamList(function), retType)
	fmt.Fprintf(file, "FujiValue %s(int argCount, FujiValue* args) {\n", wg.wrapperSymbol(function))
	fmt.Fprintf(file, "    if (argCount < %d) return NULL_VAL;\n", len(function.Parameters))
	for i, param := range function.Parameters {
		valExpr := fmt.Sprintf("args[%d]", i)
		fmt.Fprintf(file, "    %s arg%d = %s;\n", wg.cTypeForParam(param.Type), i, wg.cArgExprFuji(api, param.Type, valExpr))
	}
	if retType == "void" {
		fmt.Fprintf(file, "    %s(", function.Name)
		for i := range function.Parameters {
			if i > 0 {
				fmt.Fprintf(file, ", ")
			}
			fmt.Fprintf(file, "arg%d", i)
		}
		fmt.Fprintf(file, ");\n")
		fmt.Fprintf(file, "    return NULL_VAL;\n")
		fmt.Fprintf(file, "}\n\n")
		return
	}
	retVarT := wg.cVarTypeForReturn(retType)
	fmt.Fprintf(file, "    %s result = %s(", retVarT, function.Name)
	for i := range function.Parameters {
		if i > 0 {
			fmt.Fprintf(file, ", ")
		}
		fmt.Fprintf(file, "arg%d", i)
	}
	fmt.Fprintf(file, ");\n")
	fmt.Fprintf(file, "    %s\n", wg.cReturnExprFuji(retType, api))
	fmt.Fprintf(file, "}\n\n")
}

func (wg *WrapperGenerator) cParamList(function Function) string {
	parts := make([]string, 0, len(function.Parameters))
	for i, param := range function.Parameters {
		name := strings.TrimSpace(param.Name)
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		parts = append(parts, strings.TrimSpace(param.Type)+" "+name)
	}
	return strings.Join(parts, ", ")
}

func (wg *WrapperGenerator) cTypeForParam(t string) string {
	t = strings.TrimSpace(t)
	if strings.Contains(t, "unsigned") && strings.Contains(t, "char") && strings.Contains(t, "*") {
		return "unsigned char*"
	}
	if strings.Contains(t, "char") && strings.Contains(t, "*") {
		return "const char*"
	}
	if strings.Contains(t, "float") {
		return "float"
	}
	if strings.Contains(t, "double") {
		return "double"
	}
	if strings.Contains(t, "bool") {
		return "int"
	}
	if strings.Contains(t, "*") {
		return "void*"
	}
	// Default to the type name itself, assuming it's a known struct or typedef
	return t
}

func (wg *WrapperGenerator) cTypeForReturn(t string) string {
	t = strings.TrimSpace(t)
	if t == "void" {
		return "void"
	}
	if strings.Contains(t, "bool") {
		return "bool"
	}
	return wg.cTypeForParam(t)
}

func (wg *WrapperGenerator) isPointerType(api *API, t string) bool {
	t = strings.TrimSpace(t)
	if strings.Contains(t, "*") {
		return true
	}
	for _, td := range api.Typedefs {
		if td.Name == t {
			return wg.isPointerType(api, td.TargetType)
		}
	}
	return false
}

func isFunctionPointerParamType(t string) bool {
	t = strings.TrimSpace(t)
	if strings.Contains(t, "(*)(") {
		return true
	}
	if strings.HasSuffix(t, "Callback") {
		return true
	}
	return false
}

func (wg *WrapperGenerator) cArgExprFuji(api *API, t string, valExpr string) string {
	t = strings.TrimSpace(t)
	if isFunctionPointerParamType(t) {
		return "(" + t + ")0"
	}
	if strings.Contains(t, "unsigned") && strings.Contains(t, "char") && strings.Contains(t, "*") {
		return fmt.Sprintf(`((unsigned char*)(void*)(IS_OBJ(%s) && AS_OBJ(%s)->type == OBJ_STRING ? ((ObjString*)AS_OBJ(%s))->chars : (const char*)""))`, valExpr, valExpr, valExpr)
	}
	if strings.Contains(t, "char") && strings.Contains(t, "*") {
		return fmt.Sprintf(`(IS_OBJ(%s) && AS_OBJ(%s)->type == OBJ_STRING ? ((ObjString*)AS_OBJ(%s))->chars : "")`, valExpr, valExpr, valExpr)
	}
	if strings.Contains(t, "bool") {
		return fmt.Sprintf("(IS_BOOL(%s) ? (AS_BOOL(%s) ? 1 : 0) : 0)", valExpr, valExpr)
	}

	st := wg.findStruct(api, t)
	if st != nil {
		// Generate recursive field extraction
		fields := ""
		for i, field := range st.Fields {
			if i > 0 {
				fields += ", "
			}
			fieldValExpr := fmt.Sprintf("fuji_get_index(%s, fuji_copy_string(\"%s\", %d))", valExpr, field.Name, len(field.Name))
			fields += wg.cArgExprFuji(api, field.Type, fieldValExpr)
		}
		return fmt.Sprintf("(%s){ %s }", st.Name, fields)
	}

	if wg.isPointerType(api, t) {
		return "NULL"
	}
	if strings.Contains(t, "float") || strings.Contains(t, "double") {
		return fmt.Sprintf("(IS_NUMBER(%s) ? (double)AS_NUMBER(%s) : 0.0)", valExpr, valExpr)
	}
	if strings.Contains(t, "int") || strings.Contains(t, "long") || strings.Contains(t, "short") || strings.Contains(t, "size_t") {
		return fmt.Sprintf("(IS_NUMBER(%s) ? (int)AS_NUMBER(%s) : 0)", valExpr, valExpr)
	}
	return fmt.Sprintf("(IS_NUMBER(%s) ? (int)AS_NUMBER(%s) : 0)", valExpr, valExpr)
}

func (wg *WrapperGenerator) cVarTypeForReturn(t string) string {
	t = strings.TrimSpace(t)
	if strings.Contains(t, "bool") {
		return "bool"
	}
	if strings.Contains(t, "char") && strings.Contains(t, "*") {
		return "const char*"
	}
	if strings.Contains(t, "float") && !strings.Contains(t, "double") {
		return "float"
	}
	if strings.Contains(t, "double") {
		return "double"
	}
	if strings.Contains(t, "*") {
		return "void*"
	}
	return wg.cTypeForParam(t)
}

func (wg *WrapperGenerator) cReturnExprFuji(t string, api *API) string {
	t = strings.TrimSpace(t)
	if isFunctionPointerParamType(t) {
		return "return NULL_VAL;"
	}
	if strings.Contains(t, "char") && strings.Contains(t, "*") {
		return "return result ? fuji_copy_string(result, (int)strlen(result)) : NULL_VAL;"
	}
	if strings.Contains(t, "bool") {
		return "return BOOL_VAL(result);"
	}

	st := wg.findStruct(api, t)
	if st != nil {
		return fmt.Sprintf("return %s;", wg.cReturnStructExpr(api, st, "result"))
	}

	if strings.Contains(t, "*") {
		return "return NULL_VAL;"
	}
	return "return NUMBER_VAL((double)result);"
}

func (wg *WrapperGenerator) cReturnStructExpr(api *API, st *Struct, valExpr string) string {
	// Create a new table object and fill its fields (fuji_allocate_object returns Value).
	res := fmt.Sprintf("({ FujiValue _obj = fuji_allocate_object(%d); ", len(st.Fields))
	for _, field := range st.Fields {
		fType := strings.TrimSpace(field.Type)
		fValExpr := fmt.Sprintf("%s.%s", valExpr, field.Name)

		// Handle nested structs recursively
		nestedSt := wg.findStruct(api, fType)

		var fieldValCode string
		if nestedSt != nil {
			fieldValCode = wg.cReturnStructExpr(api, nestedSt, fValExpr)
		} else if strings.Contains(fType, "char") && strings.Contains(fType, "*") {
			fieldValCode = fmt.Sprintf("%s ? fuji_copy_string(%s, (int)strlen(%s)) : NULL_VAL", fValExpr, fValExpr, fValExpr)
		} else if strings.Contains(fType, "*") {
			fieldValCode = "NULL_VAL"
		} else if strings.Contains(fType, "bool") {
			fieldValCode = fmt.Sprintf("BOOL_VAL(%s)", fValExpr)
		} else if strings.Contains(fType, "float") || strings.Contains(fType, "double") || strings.Contains(fType, "int") || strings.Contains(fType, "long") || strings.Contains(fType, "short") || strings.Contains(fType, "size_t") || strings.Contains(fType, "unsigned") {
			fieldValCode = fmt.Sprintf("NUMBER_VAL((double)%s)", fValExpr)
		} else {
			fieldValCode = "NULL_VAL"
		}

		res += fmt.Sprintf("fuji_object_set(_obj, fuji_copy_string(\"%s\", %d), %s); ", field.Name, len(field.Name), fieldValCode)
	}
	res += "_obj; })"
	return res
}

func (wg *WrapperGenerator) exampleArgs(function Function) string {
	args := make([]string, 0, len(function.Parameters))
	for i, param := range function.Parameters {
		name := strings.TrimSpace(param.Name)
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		switch {
		case strings.Contains(param.Type, "char") && strings.Contains(param.Type, "*"):
			args = append(args, "\"text\"")
		case strings.Contains(param.Type, "*"):
			args = append(args, "0")
		default:
			args = append(args, name)
		}
	}
	return strings.Join(args, ", ")
}

func (wg *WrapperGenerator) writeFujiStructDefinition(file *os.File, structDef Struct) {
	fmt.Fprintf(file, "// Struct: %s\n", structDef.Name)
	for _, field := range structDef.Fields {
		fmt.Fprintf(file, "//   %s: %s\n", field.Name, field.Type)
	}
	fmt.Fprintf(file, "\n")
}

func (wg *WrapperGenerator) writeFujiEnumDefinition(file *os.File, enum Enum) {
	fmt.Fprintf(file, "// Enum: %s\n", enum.Name)
	for _, value := range enum.Values {
		fmt.Fprintf(file, "let %s_%s = %d;\n", enum.Name, value.Name, value.Value)
	}
	fmt.Fprintf(file, "\n")
}

func (wg *WrapperGenerator) writeFujiConstant(file *os.File, constant Constant) {
	fmt.Fprintf(file, "let %s = %s; // %s\n", constant.Name, constant.Value, constant.Type)
}

// analyzeDependencies analyzes dependencies and relationships in the API
func (wg *WrapperGenerator) analyzeDependencies(api *API) {
	// Build type dependency graph
	typeMap := make(map[string]bool)

	// Collect all defined types
	for _, structDef := range api.Structs {
		typeMap[structDef.Name] = true
	}
	for _, enum := range api.Enums {
		typeMap[enum.Name] = true
	}
	for _, typedef := range api.Typedefs {
		typeMap[typedef.Name] = true
	}

	// Analyze function parameter and return types for dependencies
	for _, function := range api.Functions {
		// Check return type
		if typeMap[function.ReturnType] {
			api.Dependencies = append(api.Dependencies, function.ReturnType)
		}

		// Check parameter types
		for _, param := range function.Parameters {
			if typeMap[param.Type] {
				api.Dependencies = append(api.Dependencies, param.Type)
			}
		}
	}

	// Remove duplicates
	uniqueDeps := make(map[string]bool)
	var cleanDeps []string
	for _, dep := range api.Dependencies {
		if !uniqueDeps[dep] {
			uniqueDeps[dep] = true
			cleanDeps = append(cleanDeps, dep)
		}
	}
	api.Dependencies = cleanDeps
}

// GenerateDocumentation generates comprehensive documentation for the API
func (wg *WrapperGenerator) GenerateDocumentation(api *API) error {
	// Generate main README
	if err := wg.generateMainReadme(api); err != nil {
		return err
	}

	// Generate API reference
	if err := wg.generateAPIReference(api); err != nil {
		return err
	}

	// Generate usage examples
	if err := wg.generateUsageExamples(api); err != nil {
		return err
	}

	// Generate HTML documentation
	if err := wg.generateHTMLDocumentation(api); err != nil {
		return err
	}

	if err := wg.generateNativeLinkingGuide(api); err != nil {
		return err
	}

	return nil
}

func (wg *WrapperGenerator) generateMainReadme(api *API) error {
	docFile := filepath.Join(wg.config.OutputDir, "README.md")
	file, err := os.Create(docFile)
	if err != nil {
		return err
	}
	defer file.Close()

	p := func(f string, a ...interface{}) { fmt.Fprintf(file, f, a...) }

	p("# %s — Fuji wrapper\n\n", api.Name)
	p("Readable **`.fuji` bindings** and a **C glue** layer for the native `%s` API, generated by **wrapgen** from your headers.\n\n", api.Name)
	p("---\n\n")

	p("## Contents\n\n")
	p("- [What you get](#what-you-get)\n")
	p("- [How it fits together](#how-it-fits-together)\n")
	p("- [Quick start](#quick-start)\n")
	p("- [Documentation map](#documentation-map)\n")
	p("- [Troubleshooting](#troubleshooting)\n\n")
	p("---\n\n")

	p("## What you get\n\n")
	p("| Output | Purpose |\n")
	p("|--------|--------|\n")
	p("| **`%s.fuji`** | Import in Fuji; declares native symbols (`// fuji:extern …`) backed by `wrapper.c`. |\n", api.Name)
	p("| **`wrapper.c`** | Bridges Fuji `Value` calls to real C functions (`fuji_wrap_%s_*`). Linked via `FUJI_NATIVE_SOURCES`. |\n", api.Name)
	p("| **`api_reference.md`** | One section per function / struct / enum / macro (large but grep-friendly). |\n")
	p("| **`examples.md`** | Copy-paste snippets for common calls. |\n")
	p("| **`docs/index.html`** | Browsable overview and signatures (open in a browser). |\n")
	p("| **`NATIVE_LINKING.md`** | How to set `FUJI_*` env vars and compile with `fuji build`. |\n")
	if wg.config.BuildSystem {
		p("| **`Makefile`**, **`CMakeLists.txt`**, **`.pc`**, **`build.sh`** | Optional templates for building the *upstream* C library (not required for `fuji build`). |\n")
	}
	p("\n---\n\n")

	p("## How it fits together\n\n")
	p("```\n")
	p("  your.fuji  ──include──▶  %s.fuji   (Fuji surface)\n", api.Name)
	p("       │\n")
	p("       └── fuji build ──▶  LLVM IR + llc + clang\n")
	p("                              │\n")
	p("                              ├──▶  libfuji_runtime.a (Fuji runtime)\n")
	p("                              └──▶  wrapper.c + lib%s (C library)\n", api.Name)
	p("```\n\n")
	p("You ship **`wrapper.c`** next to your game sources (or under `wrappers/%s/`) and point **`FUJI_NATIVE_SOURCES`** at it. The C compiler must see **`fuji_wrapgen_abi.h`** from the Fuji repo (`runtime/src/`) when compiling `wrapper.c`; `fuji build` adds that include path automatically.\n\n", api.Name)
	p("---\n\n")

	p("## Library summary\n\n")
	if len(api.Functions) > 0 {
		p("- **%d** functions\n", len(api.Functions))
	}
	if len(api.Structs) > 0 {
		p("- **%d** structs\n", len(api.Structs))
	}
	if len(api.Enums) > 0 {
		p("- **%d** enums\n", len(api.Enums))
	}
	if len(api.Constants) > 0 {
		p("- **%d** constants\n", len(api.Constants))
	}
	p("\n---\n\n")

	p("## Quick start\n\n")
	p("**Step 1.** Include the bindings (from `FUJI_WRAPPERS` / `wrappers/` layout, use `@`):\n\n")
	p("```fuji\n")
	p("#include \"@%s\"\n", api.Name)
	p("```\n\n")
	if len(api.Functions) > 0 {
		fn := api.Functions[0]
		p("**Step 2.** Call functions directly by name:\n\n")
		p("```fuji\n")
		p("let result = %s(%s);\n", fn.Name, wg.exampleArgs(fn))
		p("print(result);\n")
		p("```\n\n")
	}
	p("**Step 3.** Build or bundle:\n\n")
	p("```powershell\n")
	p("set FUJI_NATIVE_SOURCES=%s\\wrapper.c\n", wg.config.OutputDir)
	p("set FUJI_LINKFLAGS=-I<include-dir> -L<lib-dir> -l%s\n", api.Name)
	p("fuji build  mygame.fuji -o mygame.exe\n")
	p("fuji bundle mygame.fuji -o dist\\mygame\n")
	p("```\n\n")
	p("---\n\n")

	p("## Documentation map\n\n")
	p("| Doc | Use when |\n")
	p("|-----|----------|\n")
	p("| [NATIVE_LINKING.md](NATIVE_LINKING.md) | First-time native link setup |\n")
	p("| [api_reference.md](api_reference.md) | Looking up one symbol |\n")
	p("| [examples.md](examples.md) | Want pasteable `.fuji` |\n")
	p("| [docs/index.html](docs/index.html) | Prefer a browser |\n\n")
	p("---\n\n")

	p("## Troubleshooting\n\n")
	p("**Undefined symbol**  \n")
	p("Make sure `FUJI_NATIVE_SOURCES` points to `wrapper.c`.\n\n")
	p("**Missing header or library**  \n")
	p("Add `-I<dir>` for headers and `-L<dir> -l%s` for the library in `FUJI_LINKFLAGS`.\n\n", api.Name)
	p("**Unexpected return values**  \n")
	p("Check the type conversions in `wrapper.c`. Pointer and struct types may need manual adjustment for complex cases.\n\n")
	p("---\n\n")

	p("## See also\n\n")
	p("- [API Reference](api_reference.md)\n")
	p("- [Examples](examples.md)\n")
	p("- [Native linking](NATIVE_LINKING.md)\n")

	return nil
}

func (wg *WrapperGenerator) generateNativeLinkingGuide(api *API) error {
	path := filepath.Join(wg.config.OutputDir, "NATIVE_LINKING.md")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	p := func(format string, a ...interface{}) { fmt.Fprintf(f, format, a...) }
	p("# Linking `%s` with Fuji native builds\n\n", api.Name)
	p("`fuji build` compiles your `.fuji` to LLVM IR, lowers with **llc**, then runs **clang** with:\n\n")
	p("- `runtime/libfuji_runtime.a` — Fuji value model, GC hooks, `fuji_*` helpers\n")
	p("- **`wrapper.c`** (this folder) — your `%s` C API, included via `FUJI_NATIVE_SOURCES`\n", api.Name)
	p("- **Your library** — headers (`-I…`) and link line (`-L… -l%s` or a full path to `.a` / `.lib`)\n\n", api.Name)
	p("## Required environment\n\n")
	p("| Variable | Meaning |\n")
	p("|----------|--------|\n")
	p("| `FUJI_NATIVE_SOURCES` | Space-separated C sources; must include `wrapper.c` (absolute or relative path). |\n")
	p("| `FUJI_LINKFLAGS` | Extra clang arguments: `-I`, `-L`, `-l`, frameworks, etc. |\n")
	p("| `FUJI_WRAPPERS` | Optional: parent directory so `#include \"@%s\"` resolves to this folder. |\n\n", api.Name)
	p("## Header paths\n\n")
	p("Clang must find the **vendor C headers** (e.g. `%s.h`) the same names `wrapper.c` includes.\n", api.Name)
	p("Add `-I/path/to/headers` inside `FUJI_LINKFLAGS`.\n\n")
	p("Fuji also passes `-I` to `runtime/src/` so `wrapper.c` can include `fuji_wrapgen_abi.h`.\n\n")
	p("## Verify your toolchain\n\n")
	p("```text\n")
	p("fuji doctor\n")
	p("```\n\n")
	return nil
}

func (wg *WrapperGenerator) generateAPIReference(api *API) error {
	docFile := filepath.Join(wg.config.OutputDir, "api_reference.md")
	file, err := os.Create(docFile)
	if err != nil {
		return err
	}
	defer file.Close()

	p := func(f string, a ...interface{}) { fmt.Fprintf(file, f, a...) }

	p("# %s - API Reference\n\n", api.Name)

	hasTOC := len(api.Functions)+len(api.Structs)+len(api.Enums)+len(api.Constants)+len(api.Macros) > 0
	if hasTOC {
		p("## Contents\n\n")
		if len(api.Functions) > 0 {
			p("- [Functions](#functions)\n")
		}
		if len(api.Structs) > 0 {
			p("- [Structs](#structs)\n")
		}
		if len(api.Enums) > 0 {
			p("- [Enums](#enums)\n")
		}
		if len(api.Constants) > 0 {
			p("- [Constants](#constants)\n")
		}
		if len(api.Macros) > 0 {
			p("- [Macros](#macros)\n")
		}
		p("\n---\n\n")
	}

	if len(api.Functions) > 0 {
		p("## Functions\n\n")
		for _, fn := range api.Functions {
			p("### %s\n\n", fn.Name)

			params := ""
			for i, param := range fn.Parameters {
				if i > 0 {
					params += ", "
				}
				params += strings.TrimSpace(param.Type) + " " + strings.TrimSpace(param.Name)
			}
			if fn.Variadic {
				params += ", ..."
			}
			p("```c\n%s %s(%s)\n```\n\n", strings.TrimSpace(fn.ReturnType), fn.Name, params)

			fujiParams := ""
			for i, param := range fn.Parameters {
				if i > 0 {
					fujiParams += ", "
				}
				name := strings.TrimSpace(param.Name)
				if name == "" {
					name = fmt.Sprintf("arg%d", i)
				}
				fujiParams += name
			}

			p("**Fuji usage**\n\n")
			p("```fuji\n")
			if strings.TrimSpace(fn.ReturnType) != "void" {
				p("let result = %s(%s);\n", fn.Name, fujiParams)
			} else {
				p("%s(%s);\n", fn.Name, fujiParams)
			}
			p("```\n\n")

			if len(fn.Parameters) > 0 {
				p("| Parameter | Type |\n")
				p("|-----------|------|\n")
				for i, param := range fn.Parameters {
					name := strings.TrimSpace(param.Name)
					if name == "" {
						name = fmt.Sprintf("arg%d", i)
					}
					p("| `%s` | `%s` |\n", name, strings.TrimSpace(param.Type))
				}
				p("\n")
			}

			p("---\n\n")
		}
	}

	if len(api.Structs) > 0 {
		p("## Structs\n\n")
		for _, st := range api.Structs {
			p("### %s\n\n", st.Name)
			if len(st.Fields) > 0 {
				p("| Field | Type |\n")
				p("|-------|------|\n")
				for _, field := range st.Fields {
					p("| `%s` | `%s` |\n", field.Name, field.Type)
				}
				p("\n")
			}
			p("---\n\n")
		}
	}

	if len(api.Enums) > 0 {
		p("## Enums\n\n")
		for _, en := range api.Enums {
			p("### %s\n\n", en.Name)
			p("| Name | Value |\n")
			p("|------|-------|\n")
			for _, v := range en.Values {
				p("| `%s` | `%d` |\n", v.Name, v.Value)
			}
			p("\n---\n\n")
		}
	}

	if len(api.Constants) > 0 {
		p("## Constants\n\n")
		p("| Name | Type | Value |\n")
		p("|------|------|-------|\n")
		for _, c := range api.Constants {
			p("| `%s` | `%s` | `%s` |\n", c.Name, c.Type, c.Value)
		}
		p("\n")
	}

	if len(api.Macros) > 0 {
		p("## Macros\n\n")
		for _, m := range api.Macros {
			p("### %s\n\n", m.Name)
			def := m.Name
			if len(m.Parameters) > 0 {
				def += "(" + strings.Join(m.Parameters, ", ") + ")"
			}
			p("```c\n#define %s %s\n```\n\n", def, m.Value)
			p("---\n\n")
		}
	}

	return nil
}

func (wg *WrapperGenerator) generateUsageExamples(api *API) error {
	examplesFile := filepath.Join(wg.config.OutputDir, "examples.md")
	file, err := os.Create(examplesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	p := func(f string, a ...interface{}) { fmt.Fprintf(file, f, a...) }

	p("# %s - Examples\n\n", api.Name)
	p("Copy any snippet into your Fuji program and adjust the arguments.\n\n")
	p("---\n\n")

	p("## Include the library\n\n")
	p("```fuji\n")
	p("#include \"%s.fuji\"\n", api.Name)
	p("```\n\n")
	p("---\n\n")

	if len(api.Functions) > 0 {
		p("## Functions\n\n")
		limit := len(api.Functions)
		if limit > 8 {
			limit = 8
		}
		for _, fn := range api.Functions[:limit] {
			p("### %s\n\n", fn.Name)
			fujiParams := ""
			for i, param := range fn.Parameters {
				if i > 0 {
					fujiParams += ", "
				}
				name := strings.TrimSpace(param.Name)
				if name == "" {
					name = fmt.Sprintf("arg%d", i)
				}
				fujiParams += name
			}
			p("```fuji\n")
			if strings.TrimSpace(fn.ReturnType) != "void" {
				p("let result = %s(%s);\n", fn.Name, fujiParams)
				p("print(result);\n")
			} else {
				p("%s(%s);\n", fn.Name, fujiParams)
			}
			p("```\n\n")
		}
		if len(api.Functions) > 8 {
			p("*See [api_reference.md](api_reference.md) for all %d functions.*\n\n", len(api.Functions))
		}
		p("---\n\n")
	}

	if len(api.Structs) > 0 {
		p("## Structs\n\n")
		st := api.Structs[0]
		p("Structs are passed as Fuji objects with matching field names:\n\n")
		p("```fuji\n")
		p("let obj = {")
		for i, field := range st.Fields {
			if i > 0 {
				p(", ")
			}
			p(" %s: 0", field.Name)
		}
		p(" };\n")
		p("```\n\n")
		p("---\n\n")
	}

	if len(api.Enums) > 0 {
		p("## Enum values\n\n")
		en := api.Enums[0]
		p("```fuji\n")
		limit := len(en.Values)
		if limit > 4 {
			limit = 4
		}
		for _, v := range en.Values[:limit] {
			p("let %s_%s = %d;\n", en.Name, v.Name, v.Value)
		}
		p("```\n\n")
	}

	return nil
}

func (wg *WrapperGenerator) generateHTMLDocumentation(api *API) error {
	docsDir := filepath.Join(wg.config.OutputDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}

	indexFile := filepath.Join(docsDir, "index.html")
	file, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	defer file.Close()

	esc := html.EscapeString
	name := esc(api.Name)

	fmt.Fprintf(file, "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n")
	fmt.Fprintf(file, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	fmt.Fprintf(file, "<title>%s — Fuji wrapgen</title>\n", name)
	fmt.Fprintf(file, "<style>\n:root{--bg:#0f1419;--panel:#1a2332;--text:#e7ecf3;--muted:#9aa8bc;--accent:#5ec8ff;--border:#2a3648;}\n")
	fmt.Fprintf(file, "*{box-sizing:border-box;}\n")
	fmt.Fprintf(file, "body{margin:0;font:16px/1.55 system-ui,Segoe UI,Roboto,sans-serif;background:var(--bg);color:var(--text);}\n")
	fmt.Fprintf(file, "header{padding:1.75rem 2rem;background:linear-gradient(135deg,var(--panel),#141c28);border-bottom:1px solid var(--border);}\n")
	fmt.Fprintf(file, "header h1{margin:0;font-size:1.65rem;font-weight:650;}\n")
	fmt.Fprintf(file, "header p{margin:0.5rem 0 0;color:var(--muted);max-width:52rem;}\n")
	fmt.Fprintf(file, "nav a{color:var(--accent);text-decoration:none;margin-right:1.25rem;font-weight:500;}\n")
	fmt.Fprintf(file, "nav a:hover{text-decoration:underline;}\n")
	fmt.Fprintf(file, "main{max-width:56rem;margin:0 auto;padding:2rem;}\n")
	fmt.Fprintf(file, "h2{font-size:1.15rem;margin:2rem 0 0.75rem;color:var(--accent);border-bottom:1px solid var(--border);padding-bottom:0.35rem;}\n")
	fmt.Fprintf(file, ".stats{display:grid;grid-template-columns:repeat(auto-fill,minmax(9rem,1fr));gap:0.75rem;margin:1rem 0;}\n")
	fmt.Fprintf(file, ".stat{background:var(--panel);border:1px solid var(--border);border-radius:8px;padding:0.75rem 1rem;}\n")
	fmt.Fprintf(file, ".stat strong{display:block;font-size:1.35rem;color:var(--accent);}\n")
	fmt.Fprintf(file, ".stat span{font-size:0.8rem;color:var(--muted);}\n")
	fmt.Fprintf(file, ".fn{margin:1rem 0;padding:1rem 1.1rem;background:var(--panel);border:1px solid var(--border);border-radius:8px;}\n")
	fmt.Fprintf(file, ".fn h3{margin:0 0 0.5rem;font-size:1rem;font-weight:600;}\n")
	fmt.Fprintf(file, "code, .sig{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:0.82rem;background:#0d1117;color:#d6deeb;padding:0.2rem 0.45rem;border-radius:4px;word-break:break-word;}\n")
	fmt.Fprintf(file, ".sig{display:block;margin-top:0.35rem;line-height:1.45;white-space:pre-wrap;}\n")
	fmt.Fprintf(file, ".note{color:var(--muted);font-size:0.9rem;margin-top:2rem;padding-top:1rem;border-top:1px solid var(--border);}\n")
	fmt.Fprintf(file, "</style>\n</head>\n<body>\n")

	fmt.Fprintf(file, "<header><h1>%s</h1>\n", name)
	fmt.Fprintf(file, "<p>Fuji wrapgen output: C signatures parsed from your headers. Use <code>%s.fuji</code> in source and <code>../README.md</code> for the workflow.</p>\n", name)
	fmt.Fprintf(file, "<nav><a href=\"../README.md\">README</a><a href=\"../NATIVE_LINKING.md\">Native linking</a><a href=\"../api_reference.md\">Full API (Markdown)</a><a href=\"../examples.md\">Examples</a></nav>\n")
	fmt.Fprintf(file, "</header>\n<main>\n")

	fmt.Fprintf(file, "<h2>Overview</h2>\n<div class=\"stats\">\n")
	fmt.Fprintf(file, "<div class=\"stat\"><strong>%d</strong><span>functions</span></div>\n", len(api.Functions))
	fmt.Fprintf(file, "<div class=\"stat\"><strong>%d</strong><span>structs</span></div>\n", len(api.Structs))
	fmt.Fprintf(file, "<div class=\"stat\"><strong>%d</strong><span>enums</span></div>\n", len(api.Enums))
	fmt.Fprintf(file, "<div class=\"stat\"><strong>%d</strong><span>macros</span></div>\n", len(api.Macros))
	fmt.Fprintf(file, "<div class=\"stat\"><strong>%d</strong><span>typedefs</span></div>\n", len(api.Typedefs))
	fmt.Fprintf(file, "<div class=\"stat\"><strong>%d</strong><span>constants</span></div>\n", len(api.Constants))
	fmt.Fprintf(file, "</div>\n")

	const htmlFnLimit = 120
	if len(api.Functions) > 0 {
		show := len(api.Functions)
		if show > htmlFnLimit {
			show = htmlFnLimit
		}
		fmt.Fprintf(file, "<h2>Functions (first %d)</h2>\n", show)
		n := show
		for i := 0; i < n; i++ {
			fn := api.Functions[i]
			fmt.Fprintf(file, "<article class=\"fn\" id=\"fn-%s\"><h3>%s</h3>\n", esc(fn.Name), esc(fn.Name))
			var b strings.Builder
			b.WriteString(strings.TrimSpace(fn.ReturnType))
			b.WriteString(" ")
			b.WriteString(fn.Name)
			b.WriteString("(")
			for j, param := range fn.Parameters {
				if j > 0 {
					b.WriteString(", ")
				}
				b.WriteString(strings.TrimSpace(param.Type))
				b.WriteString(" ")
				b.WriteString(strings.TrimSpace(param.Name))
			}
			b.WriteString(")")
			fmt.Fprintf(file, "<span class=\"sig\">%s</span></article>\n", esc(b.String()))
		}
		if len(api.Functions) > htmlFnLimit {
			fmt.Fprintf(file, "<p class=\"note\">Showing %d of %d functions. Open <code>api_reference.md</code> for the complete list.</p>\n", htmlFnLimit, len(api.Functions))
		}
	}

	fmt.Fprintf(file, "<p class=\"note\">Generated by Fuji <strong>wrapgen</strong>. Regenerate after upgrading headers or the tool.</p>\n")
	fmt.Fprintf(file, "</main>\n</body>\n</html>\n")

	return nil
}

// GenerateBuildSystem generates build system files (Makefile, CMakeLists.txt, etc.)
func (wg *WrapperGenerator) GenerateBuildSystem(api *API) error {
	// Generate Makefile
	makefile := filepath.Join(wg.config.OutputDir, "Makefile")
	if err := wg.generateMakefile(makefile, api); err != nil {
		return err
	}

	// Generate CMakeLists.txt
	cmakeFile := filepath.Join(wg.config.OutputDir, "CMakeLists.txt")
	if err := wg.generateCMakeFile(cmakeFile, api); err != nil {
		return err
	}

	// Generate pkg-config file
	pkgConfigFile := filepath.Join(wg.config.OutputDir, api.Name+".pc")
	if err := wg.generatePkgConfigFile(pkgConfigFile, api); err != nil {
		return err
	}

	// Generate build script
	buildScript := filepath.Join(wg.config.OutputDir, "build.sh")
	if err := wg.generateBuildScript(buildScript, api); err != nil {
		return err
	}

	return nil
}

func (wg *WrapperGenerator) generateMakefile(makefile string, api *API) error {
	file, err := os.Create(makefile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# Makefile for %s wrapper\n", api.Name)
	fmt.Fprintf(file, "# Generated by WrapGen\n\n")

	fmt.Fprintf(file, "CC = gcc\n")
	fmt.Fprintf(file, "CXX = g++\n")
	fmt.Fprintf(file, "CFLAGS = -Wall -O2 -fPIC\n")
	fmt.Fprintf(file, "CXXFLAGS = -Wall -O2 -fPIC -std=c++11\n")
	fmt.Fprintf(file, "LDFLAGS = -shared\n")
	fmt.Fprintf(file, "LIBRARY = %s\n", api.Name)
	fmt.Fprintf(file, "WRAPPER = wrapper\n")
	fmt.Fprintf(file, "VERSION = 1.0.0\n\n")

	fmt.Fprintf(file, "# Source files\n")
	fmt.Fprintf(file, "SOURCES = wrapper.c\n")
	if wg.config.ComplexCPP {
		fmt.Fprintf(file, "CPP_SOURCES = wrapper.cpp\n")
	}
	fmt.Fprintf(file, "OBJECTS = $(SOURCES:.c=.o)\n")
	if wg.config.ComplexCPP {
		fmt.Fprintf(file, "CPP_OBJECTS = $(CPP_SOURCES:.cpp=.o)\n")
		fmt.Fprintf(file, "ALL_OBJECTS = $(OBJECTS) $(CPP_OBJECTS)\n")
	} else {
		fmt.Fprintf(file, "ALL_OBJECTS = $(OBJECTS)\n")
	}
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "# Default target\n")
	fmt.Fprintf(file, "all: $(WRAPPER).so $(WRAPPER).a\n\n")

	fmt.Fprintf(file, "# Shared library\n")
	fmt.Fprintf(file, "$(WRAPPER).so: $(ALL_OBJECTS)\n")
	fmt.Fprintf(file, "\t$(CC) $(LDFLAGS) -o $@ $^ -l$(LIBRARY)\n\n")

	fmt.Fprintf(file, "# Static library\n")
	fmt.Fprintf(file, "$(WRAPPER).a: $(ALL_OBJECTS)\n")
	fmt.Fprintf(file, "\tar rcs $@ $^\n\n")

	fmt.Fprintf(file, "# Compile C source files\n")
	fmt.Fprintf(file, "%%.o: %%.c\n")
	fmt.Fprintf(file, "\t$(CC) $(CFLAGS) -c $< -o $@\n\n")

	if wg.config.ComplexCPP {
		fmt.Fprintf(file, "# Compile C++ source files\n")
		fmt.Fprintf(file, "%%.o: %%.cpp\n")
		fmt.Fprintf(file, "\t$(CXX) $(CXXFLAGS) -c $< -o $@\n\n")
	}

	fmt.Fprintf(file, "# Clean\n")
	fmt.Fprintf(file, "clean:\n")
	fmt.Fprintf(file, "\trm -f $(ALL_OBJECTS) $(WRAPPER).so $(WRAPPER).a\n\n")

	fmt.Fprintf(file, "# Install\n")
	fmt.Fprintf(file, "install: all\n")
	fmt.Fprintf(file, "\tinstall -d $(DESTDIR)/usr/lib\n")
	fmt.Fprintf(file, "\tinstall -d $(DESTDIR)/usr/include\n")
	fmt.Fprintf(file, "\tinstall -m 644 $(WRAPPER).so $(DESTDIR)/usr/lib/\n")
	fmt.Fprintf(file, "\tinstall -m 644 $(WRAPPER).a $(DESTDIR)/usr/lib/\n")
	fmt.Fprintf(file, "\tinstall -m 644 *.h $(DESTDIR)/usr/include/\n\n")

	fmt.Fprintf(file, "# Uninstall\n")
	fmt.Fprintf(file, "uninstall:\n")
	fmt.Fprintf(file, "\trm -f $(DESTDIR)/usr/lib/$(WRAPPER).so\n")
	fmt.Fprintf(file, "\trm -f $(DESTDIR)/usr/lib/$(WRAPPER).a\n")
	fmt.Fprintf(file, "\trm -f $(DESTDIR)/usr/include/%s*.h\n\n", api.Name)

	fmt.Fprintf(file, "# Test\n")
	fmt.Fprintf(file, "test: all\n")
	fmt.Fprintf(file, "\t./test_%s\n", api.Name)

	return nil
}

func (wg *WrapperGenerator) generateCMakeFile(cmakeFile string, api *API) error {
	file, err := os.Create(cmakeFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# CMakeLists.txt for %s wrapper\n", api.Name)
	fmt.Fprintf(file, "# Generated by WrapGen\n\n")
	fmt.Fprintf(file, "cmake_minimum_required(VERSION 3.10)\n")
	fmt.Fprintf(file, "project(%s_wrapper VERSION 1.0.0)\n\n", api.Name)

	fmt.Fprintf(file, "# Set C and C++ standards\n")
	fmt.Fprintf(file, "set(CMAKE_C_STANDARD 99)\n")
	fmt.Fprintf(file, "set(CMAKE_CXX_STANDARD 11)\n")
	fmt.Fprintf(file, "set(CMAKE_POSITION_INDEPENDENT_CODE ON)\n\n")

	fmt.Fprintf(file, "# Find required packages\n")
	fmt.Fprintf(file, "find_package(PkgConfig REQUIRED)\n")
	fmt.Fprintf(file, "pkg_check_modules(%s REQUIRED IMPORTED_TARGET lib%s)\n\n", api.Name, api.Name)

	fmt.Fprintf(file, "# Source files\n")
	fmt.Fprintf(file, "set(WRAPPER_SOURCES wrapper.c)\n")
	if wg.config.ComplexCPP {
		fmt.Fprintf(file, "set(WRAPPER_CPP_SOURCES wrapper.cpp)\n")
		fmt.Fprintf(file, "set(WRAPPER_SOURCES ${WRAPPER_SOURCES} ${WRAPPER_CPP_SOURCES})\n")
	}
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "# Create shared library\n")
	fmt.Fprintf(file, "add_library(%s SHARED ${WRAPPER_SOURCES})\n", api.Name)
	fmt.Fprintf(file, "target_link_libraries(%s PkgConfig::%s)\n", api.Name, api.Name)
	fmt.Fprintf(file, "target_include_directories(%s PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})\n", api.Name)
	fmt.Fprintf(file, "set_target_properties(%s PROPERTIES VERSION ${PROJECT_VERSION})\n\n", api.Name)

	fmt.Fprintf(file, "# Create static library\n")
	fmt.Fprintf(file, "add_library(%s_static STATIC ${WRAPPER_SOURCES})\n", api.Name)
	fmt.Fprintf(file, "target_link_libraries(%s_static PkgConfig::%s)\n", api.Name, api.Name)
	fmt.Fprintf(file, "target_include_directories(%s_static PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})\n", api.Name)
	fmt.Fprintf(file, "set_target_properties(%s_static PROPERTIES OUTPUT_NAME %s)\n\n", api.Name, api.Name)

	fmt.Fprintf(file, "# Installation\n")
	fmt.Fprintf(file, "install(TARGETS %s %s_static\n", api.Name, api.Name)
	fmt.Fprintf(file, "    LIBRARY DESTINATION lib\n")
	fmt.Fprintf(file, "    ARCHIVE DESTINATION lib\n")
	fmt.Fprintf(file, "    RUNTIME DESTINATION bin)\n")
	fmt.Fprintf(file, "install(FILES *.h DESTINATION include)\n\n")

	fmt.Fprintf(file, "# Generate pkg-config file\n")
	fmt.Fprintf(file, "configure_file(%s.pc.in %s.pc @ONLY)\n", api.Name, api.Name)
	fmt.Fprintf(file, "install(FILES ${CMAKE_BINARY_DIR}/%s.pc DESTINATION lib/pkgconfig)\n\n", api.Name)

	fmt.Fprintf(file, "# Testing\n")
	fmt.Fprintf(file, "enable_testing()\n")
	if wg.config.IncludeTests {
		fmt.Fprintf(file, "add_executable(test_%s test_%s.c)\n", api.Name, api.Name)
		fmt.Fprintf(file, "target_link_libraries(test_%s %s)\n", api.Name, api.Name)
		fmt.Fprintf(file, "add_test(NAME test_%s COMMAND test_%s)\n", api.Name, api.Name)
	}

	return nil
}

func (wg *WrapperGenerator) generatePkgConfigFile(pkgConfigFile string, api *API) error {
	file, err := os.Create(pkgConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# pkg-config file for %s wrapper\n", api.Name)
	fmt.Fprintf(file, "# Generated by WrapGen\n\n")
	fmt.Fprintf(file, "prefix=/usr/local\n")
	fmt.Fprintf(file, "exec_prefix=${prefix}\n")
	fmt.Fprintf(file, "libdir=${exec_prefix}/lib\n")
	fmt.Fprintf(file, "includedir=${prefix}/include\n\n")
	fmt.Fprintf(file, "Name: %s\n", api.Name)
	fmt.Fprintf(file, "Description: Fuji wrapper for %s library\n", api.Name)
	fmt.Fprintf(file, "Version: 1.0.0\n")
	fmt.Fprintf(file, "Requires: %s\n", api.Name)
	fmt.Fprintf(file, "Libs: -L${libdir} -l%s\n", api.Name)
	fmt.Fprintf(file, "Cflags: -I${includedir}\n")

	return nil
}

func (wg *WrapperGenerator) generateBuildScript(buildScript string, api *API) error {
	file, err := os.Create(buildScript)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "#!/bin/bash\n")
	fmt.Fprintf(file, "# Build script for %s wrapper\n", api.Name)
	fmt.Fprintf(file, "# Generated by WrapGen\n\n")
	fmt.Fprintf(file, "set -e\n\n")
	fmt.Fprintf(file, "echo \"Building %s wrapper...\"\n", api.Name)
	fmt.Fprintf(file, "echo \"========================\"\n\n")

	fmt.Fprintf(file, "# Check dependencies\n")
	fmt.Fprintf(file, "if ! command -v pkg-config &> /dev/null; then\n")
	fmt.Fprintf(file, "    echo \"Error: pkg-config is required but not installed.\"\n")
	fmt.Fprintf(file, "    exit 1\n")
	fmt.Fprintf(file, "fi\n\n")

	fmt.Fprintf(file, "if ! pkg-config --exists %s; then\n", api.Name)
	fmt.Fprintf(file, "    echo \"Error: %s library not found.\"\n", api.Name)
	fmt.Fprintf(file, "    echo \"Please install %s development package.\"\n", api.Name)
	fmt.Fprintf(file, "    exit 1\n")
	fmt.Fprintf(file, "fi\n\n")

	fmt.Fprintf(file, "# Get compiler flags\n")
	fmt.Fprintf(file, "CFLAGS=$(pkg-config --cflags %s)\n", api.Name)
	fmt.Fprintf(file, "LIBS=$(pkg-config --libs %s)\n", api.Name)
	fmt.Fprintf(file, "echo \"CFLAGS: $CFLAGS\"\n")
	fmt.Fprintf(file, "echo \"LIBS: $LIBS\"\n\n")

	fmt.Fprintf(file, "# Build with Make\n")
	fmt.Fprintf(file, "if command -v make &> /dev/null; then\n")
	fmt.Fprintf(file, "    echo \"Building with Make...\"\n")
	fmt.Fprintf(file, "    make clean\n")
	fmt.Fprintf(file, "    make all\n")
	fmt.Fprintf(file, "    echo \"Build completed successfully!\"\n")
	fmt.Fprintf(file, "    echo \"Libraries: wrapper.so wrapper.a\"\n")
	fmt.Fprintf(file, "else\n")
	fmt.Fprintf(file, "    echo \"Make not found, building manually...\"\n")
	fmt.Fprintf(file, "    gcc -fPIC -shared -o wrapper.so wrapper.c $LIBS\n")
	fmt.Fprintf(file, "    ar rcs wrapper.a wrapper.o\n")
	fmt.Fprintf(file, "    echo \"Manual build completed!\"\n")
	fmt.Fprintf(file, "fi\n\n")

	fmt.Fprintf(file, "# Run tests if available\n")
	if wg.config.IncludeTests {
		fmt.Fprintf(file, "if [ -f \"test_%s\" ]; then\n", api.Name)
		fmt.Fprintf(file, "    echo \"Running tests...\"\n")
		fmt.Fprintf(file, "    ./test_%s\n", api.Name)
		fmt.Fprintf(file, "    echo \"Tests completed!\"\n")
		fmt.Fprintf(file, "fi\n\n")
	}

	fmt.Fprintf(file, "echo \"Build process completed!\"\n")
	fmt.Fprintf(file, "echo \"========================\"\n")
	fmt.Fprintf(file, "echo \"To install: sudo make install\"\n")
	fmt.Fprintf(file, "echo \"To test: ./test_%s\"\n", api.Name)

	// Make the script executable
	file.Chmod(0755)

	return nil
}

// GenerateTests generates test files for the wrapper
func (wg *WrapperGenerator) GenerateTests(api *API) error {
	testFile := filepath.Join(wg.config.OutputDir, "test_"+api.Name+".fuji")

	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "// Auto-generated tests for %s wrapper\n", api.Name)
	fmt.Fprintf(file, "// Generated by WrapGen\n\n")

	fmt.Fprintf(file, "// Test basic functionality\n")
	fmt.Fprintf(file, "print(\"Testing %s wrapper...\");\n\n", api.Name)

	// Generate basic function tests
	for i, function := range api.Functions {
		if i < 5 { // Limit to first 5 functions for basic testing
			if wg.isTestableFunction(function) {
				fmt.Fprintf(file, "// Test %s\n", function.Name)
				if len(function.Parameters) == 0 {
					fmt.Fprintf(file, "// let result = %s();\n", function.Name)
				} else {
					fmt.Fprintf(file, "// let result = %s(/* parameters */);\n", function.Name)
				}
				fmt.Fprintf(file, "// print(\"%s test passed\");\n\n", function.Name)
			}
		}
	}

	fmt.Fprintf(file, "print(\"All tests completed!\");\n")

	return nil
}

func (wg *WrapperGenerator) isTestableFunction(function Function) bool {
	// Skip functions that are hard to test automatically
	skipPatterns := []string{
		"init", "cleanup", "destroy", "free", "create", "open", "close",
	}

	for _, pattern := range skipPatterns {
		if strings.Contains(strings.ToLower(function.Name), pattern) {
			return false
		}
	}

	return true
}
