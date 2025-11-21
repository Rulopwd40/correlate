package utils

import (
	"encoding/json"
	"github.com/Rulopwd40/correlate/internal/models"
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
