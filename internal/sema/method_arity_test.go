package sema

import "testing"

func TestArgvMethodArityBoundsTable(t *testing.T) {
	cases := []struct {
		name      string
		wantMin   int
		wantMax   int
		wantKnown bool
	}{
		{"trim", 0, 0, true},
		{"split", 1, 1, true},
		{"slice", 2, 2, true},
		{"reduce", 1, 2, true},
		{"concat", 0, 0, false},
		{"unknownMethod", 0, 0, false},
	}
	for _, tc := range cases {
		min, max, known := argvMethodArityBounds(tc.name)
		if known != tc.wantKnown || min != tc.wantMin || max != tc.wantMax {
			t.Fatalf("%q: got min=%d max=%d known=%v want min=%d max=%d known=%v",
				tc.name, min, max, known, tc.wantMin, tc.wantMax, tc.wantKnown)
		}
	}
}
