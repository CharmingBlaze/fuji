package main

import (
	"sort"
	"strings"
)

// prefixRule: longer prefixes must appear first for longest-prefix wins.
var prefixCategoryRules = []struct {
	prefix   string
	category string
}{
	{"WindowShould", "Window"},
	{"InitAudio", "Audio"},
	{"InitWindow", "Window"},
	{"Init", "Initialisation"},
	{"Close", "Cleanup"},
	{"Unload", "Loading resources"},
	{"Load", "Loading resources"},
	{"Draw", "Drawing"},
	{"Begin", "Frame control"},
	{"End", "Frame control"},
	{"Update", "Update"},
	{"Enable", "Configuration"},
	{"Disable", "Configuration"},
	{"Check", "Checks"},
	{"Is", "Checks"},
	{"Get", "Getters"},
	{"Set", "Setters"},
	{"Play", "Audio"},
	{"Stop", "Audio"},
	{"Pause", "Audio"},
	{"Show", "UI"},
	{"Hide", "UI"},
	{"Open", "Files"},
	{"Export", "Files"},
	{"Import", "Files"},
	{"Window", "Window"},
	{"Audio", "Audio"},
	{"Text", "Text"},
	{"Image", "Textures"},
	{"Shader", "Shaders"},
	{"Camera", "Camera"},
	{"Mouse", "Input"},
	{"Keyboard", "Input"},
	{"Gamepad", "Input"},
	{"Touch", "Input"},
	{"Font", "Text"},
	{"Sound", "Audio"},
	{"Music", "Audio"},
	{"Model", "Models"},
	{"Mesh", "Models"},
	{"Ray", "Collision"},
	{"Gen", "Generation"},
	{"Mem", "Memory"},
	{"Trace", "Debug"},
}

func inferCategory(funcName string) string {
	for _, rule := range prefixCategoryRules {
		if strings.HasPrefix(funcName, rule.prefix) {
			return rule.category
		}
	}
	return "Miscellaneous"
}

func sortCategoryKeys(cats map[string][]Function) []string {
	keys := make([]string, 0, len(cats))
	for k := range cats {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortFunctionsByName(funcs []Function) []Function {
	out := append([]Function(nil), funcs...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func groupFunctionsByCategory(funcs []Function) map[string][]Function {
	out := make(map[string][]Function)
	for _, f := range funcs {
		if strings.TrimSpace(f.Name) == "" {
			continue
		}
		c := inferCategory(f.Name)
		out[c] = append(out[c], f)
	}
	for k, v := range out {
		out[k] = sortFunctionsByName(v)
	}
	return out
}

func anchorID(category string) string {
	s := strings.ToLower(category)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "misc"
	}
	return b.String()
}
