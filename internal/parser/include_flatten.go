package parser

import (
	"fmt"
	"strings"
)

// FlattenEntryIncludes expands top-level #include directives in the entry program
// into one merged Program. Requires bundle.Modules to contain every included file
// (see findImports + loadModule in loader.go).
func FlattenEntryIncludes(bundle *ProgramBundle) error {
	if bundle == nil || bundle.Entry == nil {
		return fmt.Errorf("invalid program bundle")
	}
	if len(bundle.Modules) == 0 {
		// In-memory bundles (e.g. tests) have no file paths; nothing to flatten.
		return nil
	}
	entryPath, err := BundleEntryPath(bundle)
	if err != nil {
		return err
	}
	decls, err := expandIncludes(entryPath, bundle, make(map[string]bool))
	if err != nil {
		return err
	}
	bundle.Entry = &Program{Declarations: decls}
	return nil
}

func expandIncludes(modulePath string, bundle *ProgramBundle, stack map[string]bool) ([]Decl, error) {
	if stack[modulePath] {
		return nil, fmt.Errorf("include cycle involving %q", modulePath)
	}
	stack[modulePath] = true
	defer delete(stack, modulePath)

	prog := bundle.Modules[modulePath]
	if prog == nil {
		return nil, fmt.Errorf("missing module %q (not loaded)", modulePath)
	}

	var out []Decl
	for _, d := range prog.Declarations {
		inc, ok := d.(*IncludeDecl)
		if !ok {
			out = append(out, d)
			continue
		}
		rel := strings.Trim(inc.Path.Lexeme, `"'`)
		abs, err := ResolveImportPath(modulePath, rel)
		if err != nil {
			return nil, fmt.Errorf("%s: include %q: %w", modulePath, rel, err)
		}
		inner, err := expandIncludes(abs, bundle, stack)
		if err != nil {
			return nil, err
		}
		out = append(out, inner...)
	}
	return out, nil
}
