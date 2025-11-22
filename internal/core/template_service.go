package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Rulopwd40/correlate/internal/models"
	"github.com/Rulopwd40/correlate/internal/utils"
)

var templateCache = map[string]models.Template{}

type iTemplateService interface {
	GenerateProjectTemplate(library string, identifier string, detectData string) (models.Template, error)
	SaveTemplate(template models.Template) error
	GetProjectTemplate() (models.Template, error)
	GetTemplate(templateName string) (models.Template, error)
	TemplateExists(templateName string) bool
	ListTemplates() ([]models.Template, error)
	InterpolateAndSaveTemplate(template *models.Template, packagerDir string) error
	GetPackagerRoot(library string) (string, error)
}

type TemplateService struct {
	fs  iFileService
	tr  TemplateRepository
	env *Environment
}

func NewTemplateService(fs iFileService, env *Environment) *TemplateService {
	if env == nil {
		env = DefaultEnvironment()
	}

	var tr TemplateRepository
	if env.GetTemplateRepositoryURL() != "" {
		tr = NewRestTemplateRepository(env.GetTemplateRepositoryURL())
	}

	return &TemplateService{fs, tr, env}
}

func (t *TemplateService) GenerateProjectTemplate(library string, identifier string, detectData string) (models.Template, error) {
	log.Printf("Generating template for: %s", library)

	// Use GetTemplate which has remote fallback
	template, err := t.GetTemplate(library)
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
	dir := t.env.GetCorrelatePath()
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
	dir := t.env.GetCorrelatePath()
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

func (t *TemplateService) GetTemplate(templateName string) (models.Template, error) {
	if tmpl, ok := templateCache[templateName]; ok {
		return tmpl, nil
	}

	templatesDir := t.env.GetTemplatesPath()
	templatePath := filepath.Join(templatesDir, templateName+".json")

	templateContent, err := os.ReadFile(templatePath)
	if err != nil && t.tr != nil {
		// Try remote repository if local fails
		log.Printf("Local template not found, fetching from repository: %s", templateName)
		templateContent, err = t.tr.FetchTemplate(templateName)
		if err != nil {
			return models.Template{}, fmt.Errorf("error fetching template %s: %w", templateName, err)
		}
	} else if err != nil {
		return models.Template{}, fmt.Errorf("error reading template %s: %w", templateName, err)
	}

	var template models.Template
	err = json.Unmarshal(templateContent, &template)
	if err != nil {
		return models.Template{}, fmt.Errorf("invalid template json: %w", err)
	}

	templateCache[templateName] = template

	return template, nil
}

func (t *TemplateService) TemplateExists(templateName string) bool {
	templatesDir := t.env.GetTemplatesPath()
	templatePath := filepath.Join(templatesDir, templateName+".json")
	_, err := os.Stat(templatePath)
	if err == nil {
		log.Println("Found local template:", templatePath)
		return true
	}
	// If local doesn't exist, check if we have remote repository
	if t.tr != nil {
		log.Println("Local template not found, checking remote repository:", templateName)
		_, err := t.tr.FetchTemplate(templateName)
		if err == nil {
			log.Println("Found remote template:", templateName)
			return true
		}
	}
	log.Println("Template not found:", templateName)
	return false
}

func (t *TemplateService) ListTemplates() ([]models.Template, error) {
	templatesDir := t.env.GetTemplatesPath()

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}
	var templates []models.Template
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		templateName := entry.Name()
		if !strings.HasSuffix(templateName, ".json") {
			continue
		}
		template, err := t.GetTemplate(templateName[:len(templateName)-5])
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func (t *TemplateService) InterpolateAndSaveTemplate(template *models.Template, packagerDir string) error {
	templatePath := filepath.Join(packagerDir, "template.json")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template.json: %w", err)
	}

	text := string(content)

	if template.Variables == nil {
		template.Variables = map[string]string{}
	}
	template.Variables["projectPath"] = packagerDir

	re := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.-]+)\s*\}\}`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, m := range matches {
		placeholder := m[0]
		key := m[1]

		value, ok := template.Variables[key]
		if !ok {
			fmt.Printf("No value found for placeholder: %s\n", key)
			continue
		}

		text = strings.ReplaceAll(text, placeholder, value)
	}

	err = os.WriteFile(templatePath, []byte(text), 0644)
	if err != nil {
		return fmt.Errorf("error writing interpolated template: %w", err)
	}

	fmt.Println("Template interpolated and saved")

	return nil
}

func (t *TemplateService) GetPackagerRoot(library string) (string, error) {
	// Use GetTemplate which has remote fallback built in
	template, err := t.GetTemplate(library)
	if err != nil {
		return "", err
	}

	buildFile := template.Detect["manifest"]
	if buildFile == "" {
		return "", errors.New("no manifest field found in template -> 'detect'")
	}

	// Start searching from current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rootDir, err := FindProjectRoot(currentDir, buildFile)
	if err != nil {
		return "", err
	}

	return rootDir, nil
}

func FindProjectRoot(startDir, buildFile string) (string, error) {
	currentDir := startDir

	for {
		buildFilePath := filepath.Join(currentDir, buildFile)

		if _, err := os.Stat(buildFilePath); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}

		currentDir = parent
	}

	return "", fmt.Errorf("build file '%s' not found", buildFile)
}
