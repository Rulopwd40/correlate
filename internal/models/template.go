package models

type Template struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables"`
	Detect      map[string]string `json:"detect"`
	Steps       []Step            `json:"steps"`
}

type Step struct {
	Name      string            `json:"name"`
	Type      string            `json:"type,omitempty"`   // "command" (default) or "script"
	Cmd       string            `json:"cmd,omitempty"`    // Single command
	Script    []string          `json:"script,omitempty"` // Multi-line script
	Workdir   string            `json:"workdir"`
	Variables map[string]string `json:"variables,omitempty"`
	Outputs   map[string]string `json:"outputs,omitempty"` // Capture outputs for next steps
}
