package util

import "testing"

func TestAddQuotes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"foo.bar", "foo.bar"},
		{"foo.bar[baz]", `foo.bar["baz"]`},
		{"foo.bar[0]", "foo.bar[0]"},
		{`foo.bar["baz"]`, `foo.bar["baz"]`},
		{"foo.bar[1s]", `foo.bar["1s"]`},
		{"foo.bar[my-key]", `foo.bar["my-key"]`},
		{"foo.bar[a.b]", `foo.bar["a.b"]`},
		{"foo.bar[_]", `foo.bar["_"]`},
		{"module.a[x].foo.bar[0]", `module.a["x"].foo.bar[0]`},
		{`module.a["x"].foo.bar[y]`, `module.a["x"].foo.bar["y"]`},
		{"module.a[0].module.b[c].foo.bar[d]", `module.a[0].module.b["c"].foo.bar["d"]`},
		// a key of digits only is taken as a number, as the shell left nothing to tell it
		{"foo.bar[10]", "foo.bar[10]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := AddQuotes(tt.input); got != tt.want {
				t.Errorf("AddQuotes(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTargetAddress(t *testing.T) {
	tests := []struct {
		arg     string
		address string
		ok      bool
	}{
		{"-target=foo.bar", "foo.bar", true},
		{"--target=foo.bar[baz]", "foo.bar[baz]", true},
		{"-target", "", false},
		{"-target-file=targets.txt", "", false},
		{"foo.bar", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			address, ok := TargetAddress(tt.arg)
			if address != tt.address || ok != tt.ok {
				t.Errorf("TargetAddress(%q) = %q, %v, want %q, %v", tt.arg, address, ok, tt.address, tt.ok)
			}
		})
	}
}
