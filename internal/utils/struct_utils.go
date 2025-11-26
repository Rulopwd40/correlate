package utils

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/Rulopwd40/correlate/internal/models"
	"github.com/Rulopwd40/correlate/internal/pipeline"
)

func ParseTemplate(content []byte) (models.Template, error) {
	template := models.Template{}
	err := json.Unmarshal(content, &template)
	if err != nil {
		return template, err
	}
	return template, nil
}

func ParseReferences(content []byte) (models.References, error) {
	references := models.References{}
	err := json.Unmarshal(content, &references)
	if err != nil {
		return references, err
	}
	return references, nil
}

func ParseConfig(content []byte) (models.Config, error) {
	config := models.Config{}
	err := json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// MakeTasks creates pipeline tasks from template steps with variable resolution
// Supports both new context-based approach and legacy single directory parameter
func MakeTasks(steps []models.Step, ref models.Reference, contextOrDir interface{}) ([]pipeline.Task, error) {
	// Build context with all available variables
	variables := make(map[string]string)

	// Handle both old (string) and new (map) calling conventions
	switch v := contextOrDir.(type) {
	case string:
		// Legacy: single directory string
		variables["identifier"] = ref.Identifier
		variables["sourceDir"] = v
		variables["targetDir"] = v
		variables["dir"] = v // Keep compatibility with old $2
	case map[string]string:
		// New: full context map
		variables["identifier"] = ref.Identifier
		for key, value := range v {
			variables[key] = value
		}
		// Ensure these keys exist for backward compatibility
		if _, ok := variables["sourceDir"]; !ok {
			if dir, hasDir := v["dir"]; hasDir {
				variables["sourceDir"] = dir
			}
		}
		if _, ok := variables["targetDir"]; !ok {
			if dir, hasDir := v["dir"]; hasDir {
				variables["targetDir"] = dir
			}
		}
	}

	tasks := make([]pipeline.Task, 0, len(steps))

	for _, step := range steps {
		var cmd string

		// Merge step-level variables with global variables
		stepVariables := make(map[string]string)
		for k, v := range variables {
			stepVariables[k] = v
		}
		for k, v := range step.Variables {
			stepVariables[k] = resolveVariables(v, variables)
		}

		// Determine command type
		if step.Type == "script" && len(step.Script) > 0 {
			// Join multi-line script into single command
			cmd = strings.Join(step.Script, " && ")
		} else {
			cmd = step.Cmd
		}

		// Resolve all variables
		cmd = resolveVariables(cmd, stepVariables)
		workdir := resolveVariables(step.Workdir, stepVariables)

		tasks = append(tasks, pipeline.Task{
			Cmd:     cmd,
			Name:    "[" + ref.Identifier + "] " + step.Name,
			Workdir: workdir,
			Outputs: step.Outputs,
		})
	}

	return tasks, nil
}

func resolveVariables(template string, variables map[string]string) string {
	result := template
	re := regexp.MustCompile(`\{\{([a-zA-Z0-9_.-]+)\}\}`)

	matches := re.FindAllStringSubmatch(template, -1)
	for _, match := range matches {
		placeholder := match[0]
		key := match[1]

		if value, ok := variables[key]; ok {
			result = strings.ReplaceAll(result, placeholder, value)
		}
	}

	return result
}
