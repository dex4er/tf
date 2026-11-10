package util

import (
	"os"
	"strings"
)

// The variables that coding agents set in the environment of the commands
// that they run, and not in the terminal of a human: AI_AGENT, the convention
// shared by Claude Code, opencode and Copilot in VS Code, then the variables of
// each agent. CLAUDECODE is left out, as the IDE extensions of Claude Code set
// it in the terminals that a human uses too. TTY detection is no substitute:
// CI runs without a terminal and wants the colors.
var agentVariables = []string{
	"AI_AGENT",
	"CLAUDE_CODE_CHILD_SESSION",
	"CODEX_CI",
	"CODEX_THREAD_ID",
	"OPENCODE",
	"COPILOT_CLI",
	"COPILOT_AGENT",
}

// Reports whether tf runs under a coding agent. `TF_AGENT` set to true or false
// overrides the detection.
func IsAgent() bool {
	switch strings.ToLower(os.Getenv("TF_AGENT")) {
	case "1", "true", "yes":
		return true
	case "0", "false", "no":
		return false
	}

	for _, name := range agentVariables {
		if os.Getenv(name) != "" {
			return true
		}
	}

	return false
}

// Reports whether the colors are off by default: under a coding agent, or with
// `NO_COLOR` (https://no-color.org), which terraform and tofu do not read.
func NoColor() bool {
	return os.Getenv("NO_COLOR") != "" || IsAgent()
}
