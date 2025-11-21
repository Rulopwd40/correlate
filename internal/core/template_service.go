package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rulopwd40/correlate/internal/models"
	"github.com/Rulopwd40/correlate/internal/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type iTemplateService interface {
	GenerateProjectTemplate(library string, identifier string, detectData string) (models.Template, error)
	SaveTemplate(template models.Template) error
	GetProjectTemplate() (models.Template, error)
}

const TEMPLATES_PATH = "templates"
const PROJECT_TEMPLATES_PATH = ".correlate"

type TemplateService struct {
	fs iFileService
}

func NewTemplateService(fs iFileService) *TemplateService {
	return &TemplateService{fs}

}

func (t *TemplateService) GenerateProjectTemplate(library string, identifier string, detectData string) (models.Template, error) {
	log.Printf("Generating template for: %s", TEMPLATES_PATH+"/"+library)

	templateJson, err := t.fs.GetJSON(TEMPLATES_PATH, library)
	if err != nil {
		return models.Template{}, err
	}

	template, err := utils.ParseTemplate(templateJson)
	if err != nil {
		return models.Template{}, err
	}
	//DetectData must be an Absolute path of manifest location
	log.Printf("Getting manifest for: %s", detectData)
	manifest, err := t.fs.GetFile(detectData, template.Detect["manifest"])
	if err != nil {
		return models.Template{}, err
	}

	pattern := template.Detect["searchPattern"]
	pattern = strings.Replace(pattern, "{{identifier}}", identifier, -1)
	template.Detect["searchPattern"] = pattern

	if !bytes.Contains(manifest, []byte(pattern)) {
		return models.Template{}, errors.New("No match for " + pattern)
	}

	err = t.SaveTemplate(template)
	if err != nil {
		return models.Template{}, err
	}

	return template, nil
}

func (t *TemplateService) SaveTemplate(template models.Template) error {
	dir := CORRELATE_PATH
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(template); err != nil {
		return err
	}

	filePath := filepath.Join(dir, "template.json")
	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return err
	}

	fmt.Println("template.json generated successfully.")
	return nil
}

func (t *TemplateService) GetProjectTemplate() (models.Template, error) {
	dir := CORRELATE_PATH
	templateJson, err := t.fs.GetJSON(dir, "template")
	if err != nil {
		return models.Template{}, err
	}
	template, err := utils.ParseTemplate(templateJson)
	if err != nil {
		return models.Template{}, err
	}
	return template, nil
}
