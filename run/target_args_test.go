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
		{"-out FILE", []string{"-out", "tfplan2", "foo.bar"}, false, []string{"-out", "tfplan2", "-target=foo.bar"}},
		{"-var NAME=VALUE", []string{"-var", "x=1", "--var-file", "x.tfvars"}, false, []string{"-var", "x=1", "--var-file", "x.tfvars"}},
		{"a value option last", []string{"-out"}, false, []string{"-out"}},
		{"a boolean option before an address", []string{"-auto-approve", "foo.bar"}, false, []string{"-auto-approve", "-target=foo.bar"}},
		{"-replace=", []string{"-replace=foo.bar[baz]"}, false, []string{`-replace=foo.bar["baz"]`}},
		{"--replace ADDRESS", []string{"--replace", "foo.bar[baz]"}, false, []string{"--replace", `foo.bar["baz"]`}},
		{"-exclude=", []string{"--exclude=foo.bar[baz]"}, false, []string{`-exclude=foo.bar["baz"]`}},
		{"-exclude ADDRESS", []string{"-exclude", "foo.bar[baz]", "foo.qux"}, false, []string{"-exclude", `foo.bar["baz"]`, "-target=foo.qux"}},
		{"-invoke=", []string{"-invoke=action.foo.bar[baz]"}, false, []string{`-invoke=action.foo.bar["baz"]`}},
		{"-var= keeps its brackets", []string{"-var=x=[a]"}, false, []string{"-var=x=[a]"}},
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
		{"-plan FILE last", []string{"-plan", "tfplan", "-auto-approve"}, []string{"-auto-approve", "tfplan"}},
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

func TestSplitAddresses(t *testing.T) {
	options, addresses := splitAddresses([]string{"-lock-timeout", "10s", "foo.bar[baz]", "-allow-missing", "foo.qux"})
	if want := []string{"-lock-timeout", "10s", "-allow-missing"}; !slices.Equal(options, want) {
		t.Errorf("options = %q, want %q", options, want)
	}
	if want := []string{`foo.bar["baz"]`, "foo.qux"}; !slices.Equal(addresses, want) {
		t.Errorf("addresses = %q, want %q", addresses, want)
	}
}

func TestQuotedArgs(t *testing.T) {
	got := quotedArgs([]string{"mv", "-state", "x.tfstate", "foo.bar[a]", "foo.bar[b]"})
	if want := []string{"mv", "-state", "x.tfstate", `foo.bar["a"]`, `foo.bar["b"]`}; !slices.Equal(got, want) {
		t.Errorf("quotedArgs = %q, want %q", got, want)
	}
}

func TestImportArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		want     []string
		resource string
		noShow   bool
	}{
		{"address and ID", []string{"foo.bar[baz]", "some", "id"}, []string{`foo.bar["baz"]`, "some id"}, `foo.bar["baz"]`, false},
		{"-var NAME=VALUE", []string{"-var", "x=1", "foo.bar", "id"}, []string{"-var", "x=1", "foo.bar", "id"}, "foo.bar", false},
		{"-no-show", []string{"-no-show", "foo.bar", "id"}, []string{"foo.bar", "id"}, "foo.bar", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, resource, noShow := importArgs(tt.args)
			if !slices.Equal(got, tt.want) || resource != tt.resource || noShow != tt.noShow {
				t.Errorf("importArgs(%q) = %q, %q, %v, want %q, %q, %v", tt.args, got, resource, noShow, tt.want, tt.resource, tt.noShow)
			}
		})
	}
}
