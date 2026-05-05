package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// WrapGen generates Fuji bindings and C glue code from C/C++ header files.
type WrapGenConfig struct {
	LibraryName   string
	InputHeaders  []string
	OutputDir     string
	Language      string
	Documentation bool
	BuildSystem   bool
	IncludeTests  bool
	ComplexCPP    bool // Support C++ templates, classes
	Verbose       bool
}

func main() {
	config := parseFlags()

	if config.Verbose {
		fmt.Printf("wrapgen: generating bindings for %q\n", config.LibraryName)
		fmt.Printf("  headers : %s\n", strings.Join(config.InputHeaders, ", "))
		fmt.Printf("  output  : %s\n\n", config.OutputDir)
	}

	if err := validateConfig(config); err != nil {
		log.Fatalf("error: %v\n\nRun wrapgen -help for usage.", err)
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("error: could not create output directory: %v", err)
	}

	generator := NewWrapperGenerator(config)

	api, err := generator.ParseHeaders()
	if err != nil {
		log.Fatalf("error: failed to parse headers: %v", err)
	}

	if err := generator.GenerateWrapper(api); err != nil {
		log.Fatalf("error: failed to generate wrapper: %v", err)
	}

	if config.BuildSystem {
		if err := generator.GenerateBuildSystem(api); err != nil {
			log.Fatalf("error: failed to generate build system: %v", err)
		}
	}

	if config.Documentation {
		if err := generator.GenerateDocumentation(api); err != nil {
			log.Fatalf("error: failed to generate documentation: %v", err)
		}
	}

	if config.IncludeTests {
		if err := generator.GenerateTests(api); err != nil {
			log.Fatalf("error: failed to generate tests: %v", err)
		}
	}

	printSuccess(config, api)
}

func printSuccess(config *WrapGenConfig, api *API) {
	out := config.OutputDir
	lib := config.LibraryName

	fmt.Printf("Wrapper generated: %s\n\n", lib)
	fmt.Printf("  Output folder : %s\n", out)
	fmt.Printf("  Fuji bindings : %s/%s.fuji\n", out, lib)
	fmt.Printf("  C glue         : %s/wrapper.c\n", out)
	if config.BuildSystem {
		fmt.Printf("  Build helpers  : %s/Makefile, CMakeLists.txt, %s.pc, build.sh\n", out, lib)
	}
	if config.Documentation {
		fmt.Printf("  Docs           : %s/README.md, api_reference.md, examples.md\n", out)
		fmt.Printf("  HTML           : %s/docs/index.html\n", out)
		fmt.Printf("  Linking guide  : %s/NATIVE_LINKING.md\n", out)
	}
	fmt.Printf("\n")
	fmt.Printf("  Functions : %d\n", len(api.Functions))
	fmt.Printf("  Structs   : %d\n", len(api.Structs))
	fmt.Printf("  Enums     : %d\n", len(api.Enums))
	fmt.Printf("  Constants : %d\n", len(api.Constants))

	fmt.Printf(`
To use in your Fuji program:

  #include "@%s"

To build / bundle (set paths to your C library headers and .a/.lib):

  set FUJI_NATIVE_SOURCES=%s\wrapper.c
  set FUJI_LINKFLAGS=-I<include-dir> -L<lib-dir> -l%s
  fuji build mygame.fuji -o mygame.exe
  fuji bundle mygame.fuji -o dist\mygame

See %s/README.md and %s/NATIVE_LINKING.md for the full guide.
`, lib, out, lib, out, out)
}

func parseFlags() *WrapGenConfig {
	config := &WrapGenConfig{}

	flag.StringVar(&config.LibraryName, "name", "", "Library name (required, e.g. raylib)")
	flag.StringVar(&config.OutputDir, "out", "./wrapper", "Output directory for generated files")
	flag.StringVar(&config.Language, "lang", "fuji", "Target language: fuji (.fuji output) (.fuji output)")
	flag.BoolVar(&config.Documentation, "docs", true, "Generate documentation (README, API reference, HTML, linking guide)")
	flag.BoolVar(&config.BuildSystem, "build", true, "Generate Makefile, CMakeLists, pkg-config stub, build.sh (pass -build=false to skip)")
	flag.BoolVar(&config.IncludeTests, "tests", false, "Generate placeholder tests (optional)")
	flag.BoolVar(&config.ComplexCPP, "cpp", true, "Parse C++ features")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")

	var headers string
	flag.StringVar(&headers, "headers", "", "Comma-separated list of header files to parse")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `wrapgen - Generate readable .fuji bindings from C/C++ headers

USAGE
  wrapgen -name <lib> -headers <file.h>[,<file2.h>] [options]
  (same flags; fuji forwards: fuji wrap ...)

OPTIONS
  -name     Library name (required)
  -headers  Comma-separated header files to parse (required)
  -out      Output directory (default: ./wrapper)
  -docs     Generate documentation (default: true)
  -build    Generate Makefile, CMake, pkg-config stub, build.sh (default: true; -build=false to skip)
  -tests    Generate placeholder tests (default: false)
  -v        Verbose output

EXAMPLES
  wrapgen -name raylib -headers raylib.h -out wrappers\raylib
  wrapgen -name sqlite3 -headers sqlite3.h -out wrappers\sqlite3
  wrapgen -name mylib -headers include\mylib.h,include\mylib_ext.h -out wrappers\mylib
`)
	}

	flag.Parse()

	if headers != "" {
		for _, h := range strings.Split(headers, ",") {
			if h = strings.TrimSpace(h); h != "" {
				config.InputHeaders = append(config.InputHeaders, h)
			}
		}
	}

	return config
}

func validateConfig(config *WrapGenConfig) error {
	if config.LibraryName == "" {
		return fmt.Errorf("-name is required (e.g. -name raylib)")
	}
	if len(config.InputHeaders) == 0 {
		return fmt.Errorf("-headers is required (e.g. -headers raylib.h)")
	}
	for _, header := range config.InputHeaders {
		if _, err := os.Stat(header); os.IsNotExist(err) {
			return fmt.Errorf("header file not found: %s", header)
		}
	}
	return nil
}
