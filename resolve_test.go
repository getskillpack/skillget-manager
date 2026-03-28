package skillgetmanager

import "testing"

func TestParseNameVersion(t *testing.T) {
	tests := []struct {
		in       string
		wantName string
		wantVer  string
	}{
		{"foo", "foo", ""},
		{"foo@1.2.3", "foo", "1.2.3"},
		{"@scope/pkg@2.0.0", "@scope/pkg", "2.0.0"},
		{"nover@at", "nover", "at"},
	}
	for _, tc := range tests {
		nv := ParseNameVersion(tc.in)
		if nv.Name != tc.wantName || nv.Version != tc.wantVer {
			t.Fatalf("ParseNameVersion(%q) = %+v; want name=%q version=%q", tc.in, nv, tc.wantName, tc.wantVer)
		}
	}
}
