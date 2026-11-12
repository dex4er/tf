package run

import (
	"slices"
	"testing"

	"github.com/dex4er/tf/progress"
)

func TestParseProgressOptions(t *testing.T) {
	agent := map[string]string{"TF_AGENT": "true"}

	tests := []struct {
		name string
		env  map[string]string
		args []string
		want progressOptions
	}{
		{
			"defaults", nil,
			[]string{"-auto-approve", "tfplan"},
			progressOptions{Short, progress.Counters, false, true, []string{"-auto-approve", "tfplan"}},
		},
		{
			"TF_IN_AUTOMATION", map[string]string{"TF_IN_AUTOMATION": "1"},
			nil,
			progressOptions{Short, progress.Verbose, false, true, []string{}},
		},
		{
			"agent puts -no-color first, before the plan file", agent,
			[]string{"-auto-approve", "tfplan"},
			progressOptions{Short, progress.Quiet, true, true, []string{"-no-color", "-auto-approve", "tfplan"}},
		},
		{
			"NO_COLOR puts -no-color first and keeps the progress", map[string]string{"NO_COLOR": "1"},
			[]string{"tfplan"},
			progressOptions{Short, progress.Counters, true, true, []string{"-no-color", "tfplan"}},
		},
		{
			"agent with -no-color keeps it once", agent,
			[]string{"-auto-approve", "-no-color"},
			progressOptions{Short, progress.Quiet, true, true, []string{"-auto-approve", "-no-color"}},
		},
		{
			"agent loses to TF_PROGRESS_FORMAT", map[string]string{"TF_AGENT": "true", "TF_PROGRESS_FORMAT": progress.Fan},
			nil,
			progressOptions{Short, progress.Fan, true, true, []string{"-no-color"}},
		},
		{
			"agent loses to an option", agent,
			[]string{"-verbose"},
			progressOptions{Short, progress.Verbose, true, true, []string{"-no-color"}},
		},
		{
			"TF_PROGRESS_FORMAT loses to an option", map[string]string{"TF_PROGRESS_FORMAT": progress.Dots},
			[]string{"--fan"},
			progressOptions{Short, progress.Fan, false, true, []string{}},
		},
		{
			"TF_PLAN_FORMAT loses to an option", map[string]string{"TF_PLAN_FORMAT": Compact},
			[]string{"-full"},
			progressOptions{Full, progress.Counters, false, true, []string{}},
		},
		{
			"-no-colors", nil,
			[]string{"-no-colors"},
			progressOptions{Short, progress.Counters, true, true, []string{"-no-color"}},
		},
		{
			"-verbatim", nil,
			[]string{"-verbatim"},
			progressOptions{Full, progress.Verbatim, false, true, []string{}},
		},
		{
			"-no-outputs=false", nil,
			[]string{"--no-outputs=false"},
			progressOptions{Short, progress.Counters, false, false, []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TF_AGENT=false keeps the agent of the environment of the test out
			for _, name := range []string{"TF_IN_AUTOMATION", "TF_PLAN_FORMAT", "TF_PROGRESS_FORMAT", "NO_COLOR"} {
				t.Setenv(name, "")
			}
			t.Setenv("TF_AGENT", "false")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got := parseProgressOptions(tt.args)
			if got.planFormat != tt.want.planFormat ||
				got.progressFormat != tt.want.progressFormat ||
				got.noColor != tt.want.noColor ||
				got.noOutputs != tt.want.noOutputs ||
				!slices.Equal(got.args, tt.want.args) {
				t.Errorf("parseProgressOptions(%q) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}
