package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindDotVersionFile(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  string
	}{
		{"none", nil, ""},
		{"terraform", []string{TerraformVersionFile}, TerraformVersionFile},
		{"opentofu", []string{OpentofuVersionFile}, OpentofuVersionFile},
		{"opentofu wins in the same directory", []string{TerraformVersionFile, OpentofuVersionFile}, OpentofuVersionFile},
		{"from a parent", []string{"../" + OpentofuVersionFile}, OpentofuVersionFile},
		{"from the parent of a parent", []string{"../../" + TerraformVersionFile}, TerraformVersionFile},
		{"the nearest wins", []string{"../" + OpentofuVersionFile, TerraformVersionFile}, TerraformVersionFile},
		{"the nearest wins over opentofu", []string{"../../" + OpentofuVersionFile, "../" + TerraformVersionFile}, TerraformVersionFile},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// two levels below the temporary directory, so that a file can be put
			// in a parent without leaving it
			dir := filepath.Join(t.TempDir(), "a", "b")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			for _, file := range tt.files {
				if err := os.WriteFile(filepath.Join(dir, file), nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			t.Chdir(dir)

			if got := FindDotVersionFile(); got != tt.want {
				t.Errorf("FindDotVersionFile() = %q, want %q", got, tt.want)
			}
		})
	}
}
