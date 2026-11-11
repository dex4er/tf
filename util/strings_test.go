package util

import "testing"

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
