package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Rulopwd40/correlate/internal/models"
	"github.com/Rulopwd40/correlate/internal/pipeline"
	"strings"
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

func MakeTasks(stepMaps []map[string]string, ref models.Reference, dir string) ([]pipeline.Task, error) {
	tasks := make([]pipeline.Task, 0, len(stepMaps))

	for _, step := range stepMaps {
		cmd, ok1 := step["cmd"]
		name, ok2 := step["name"]
		workdir, ok3 := step["workdir"]

		if !ok1 || !ok2 || !ok3 {
			return nil, fmt.Errorf("invalid step: missing cmd, name or workdir")
		}

		// Reemplazos estándar
		cmd = strings.ReplaceAll(cmd, "$1", ref.Identifier)
		cmd = strings.ReplaceAll(cmd, "$2", dir)

		workdir = strings.ReplaceAll(workdir, "$1", ref.Identifier)
		workdir = strings.ReplaceAll(workdir, "$2", dir)

		task := pipeline.Task{
			Cmd:     cmd,
			Name:    "[" + ref.Identifier + "]" + name,
			Workdir: workdir,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
