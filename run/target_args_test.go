package run

import (
	"os"
	"slices"
	"testing"
)

func TestTargetArgs(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile("tfplan", nil, 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		args  []string
		files bool
		want  []string
	}{
		{"none", nil, false, []string{}},
		{"options", []string{"-parallelism=30", "--compact"}, false, []string{"-parallelism=30", "--compact"}},
		{"address", []string{"foo.bar[baz]"}, false, []string{`-target=foo.bar["baz"]`}},
		{"addresses", []string{"foo.bar", "foo.baz[0]"}, false, []string{"-target=foo.bar", "-target=foo.baz[0]"}},
		{"-target=", []string{"-target=foo.bar[baz]"}, false, []string{`-target=foo.bar["baz"]`}},
		{"--target=", []string{"--target=foo.bar"}, false, []string{"-target=foo.bar"}},
		{"-target ADDRESS", []string{"-target", "foo.bar[baz]"}, false, []string{`-target=foo.bar["baz"]`}},
		{"--target ADDRESS", []string{"--target", "foo.bar"}, false, []string{"-target=foo.bar"}},
		{"order kept", []string{"-auto-approve", "foo.bar", "-parallelism=30"}, false, []string{"-auto-approve", "-target=foo.bar", "-parallelism=30"}},
		{"file without files", []string{"tfplan"}, false, []string{"-target=tfplan"}},
		{"file with files", []string{"tfplan"}, true, []string{"tfplan"}},
		{"no file with files", []string{"foo.bar"}, true, []string{"-target=foo.bar"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := targetArgs(tt.args, tt.files); !slices.Equal(got, tt.want) {
				t.Errorf("targetArgs(%q, %v) = %q, want %q", tt.args, tt.files, got, tt.want)
			}
		})
	}
}

func TestApplyArgs(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile("tfplan", nil, 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"none", nil, []string{}},
		{"plan file", []string{"tfplan"}, []string{"tfplan"}},
		{"-plan= last", []string{"-plan=tfplan", "-auto-approve"}, []string{"-auto-approve", "tfplan"}},
		{"--plan= last", []string{"--plan=tfplan", "-auto-approve"}, []string{"-auto-approve", "tfplan"}},
		{"-plan= of a file not there yet", []string{"-plan=missing"}, []string{"missing"}},
		{"address", []string{"-auto-approve", "foo.bar[baz]"}, []string{"-auto-approve", `-target=foo.bar["baz"]`}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyArgs(tt.args); !slices.Equal(got, tt.want) {
				t.Errorf("applyArgs(%q) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}
