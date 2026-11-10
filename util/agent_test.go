package util

import "testing"

func clearAgentEnv(t *testing.T) {
	t.Helper()
	t.Setenv("TF_AGENT", "")
	t.Setenv("NO_COLOR", "")
	for _, name := range agentVariables {
		t.Setenv(name, "")
	}
}

func TestIsAgent(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{"none", nil, false},
		{"AI_AGENT", map[string]string{"AI_AGENT": "claude-code_2-1-322_agent"}, true},
		{"Claude Code", map[string]string{"CLAUDE_CODE_CHILD_SESSION": "1"}, true},
		{"Claude Code IDE terminal", map[string]string{"CLAUDECODE": "1"}, false},
		{"Codex", map[string]string{"CODEX_CI": "1"}, true},
		{"opencode", map[string]string{"OPENCODE": "1"}, true},
		{"Copilot CLI", map[string]string{"COPILOT_CLI": "1"}, true},
		{"Copilot in VS Code", map[string]string{"COPILOT_AGENT": "1"}, true},
		{"TF_AGENT=true", map[string]string{"TF_AGENT": "true"}, true},
		{"TF_AGENT=false wins", map[string]string{"TF_AGENT": "false", "AI_AGENT": "opencode"}, false},
		{"NO_COLOR is not an agent", map[string]string{"NO_COLOR": "1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearAgentEnv(t)
			t.Setenv("CLAUDECODE", "")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			if got := IsAgent(); got != tt.want {
				t.Errorf("IsAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoColor(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{"none", nil, false},
		{"NO_COLOR", map[string]string{"NO_COLOR": "1"}, true},
		{"agent", map[string]string{"OPENCODE": "1"}, true},
		{"agent detection off", map[string]string{"TF_AGENT": "false", "OPENCODE": "1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearAgentEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			if got := NoColor(); got != tt.want {
				t.Errorf("NoColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
