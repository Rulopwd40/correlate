package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rulopwd40/correlate/internal/utils"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Rulopwd40/correlate/internal/models"
)

var templateCache = map[string]models.Template{}

type iFileService interface {
	GetPackagerRoot(library string) (string, error)
	TemplateExists(templateName string) bool
	GetJSON(path string, name string) ([]byte, error)
	GetFile(path string, name string) ([]byte, error)
	FindAllFiles(root string, name string) ([]string, error)
	FindAllFilesThatContains(root string, filename string, identifier string) ([]string, error)
}

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (fs *FileService) GetJSON(path string, name string) ([]byte, error) {
	fullPath, err := fs.resolvePath(path, name+".json")
	if err != nil {
		return nil, err
	}

	return os.ReadFile(fullPath)
}
func (fs *FileService) GetFile(path string, name string) ([]byte, error) {
	fullPath, err := fs.resolvePath(path, name)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(fullPath)
}

func (fs *FileService) FindAllFiles(root string, name string) ([]string, error) {
	var results []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // error al acceder a un archivo/directorio
		}

		if !info.IsDir() && info.Name() == name {
			results = append(results, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (fs *FileService) FindAllFilesThatContains(root string, filename string, identifier string) ([]string, error) {
	var results []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == filename {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if bytes.Contains(data, []byte(identifier)) {
				results = append(results, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (fs *FileService) resolvePath(path, name string) (string, error) {
	base, err := fs.baseDirFrom(path)
	if err != nil {
		return "", err
	}

	full := filepath.Join(base, name)
	return full, nil
}

func (fs *FileService) baseDirFrom(path string) (string, error) {
	if path == "" {
		return os.Getwd()
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, path), nil
}

func (fs *FileService) GetTemplatesDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, "templates"), nil
}

func (fs *FileService) GetTemplate(templateName string) (models.Template, error) {
	if t, ok := templateCache[templateName]; ok {
		return t, nil
	}

	templatesDir, err := fs.GetTemplatesDir()
	if err != nil {
		return models.Template{}, err
	}

	templatePath := filepath.Join(templatesDir, templateName+".json")

	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
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
func (fs *FileService) TemplateExists(templateName string) bool {
	templatesDir, err := fs.GetTemplatesDir()
	if err != nil {
		return false
	}
	log.Println("Searching template:", templateName, "in template dir:", templatesDir)
	templatePath := filepath.Join(templatesDir, templateName+".json")
	log.Println("Checking template:", templatePath)
	_, err = os.Stat(templatePath)
	return err == nil
}

func (fs *FileService) ListTemplates() ([]models.Template, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	templatesDir := filepath.Join(currentDir, "templates")

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
		template, err := fs.GetTemplate(templateName[:len(templateName)-5])
		if err != nil {

			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func (fs *FileService) InterpolateAndSaveTemplate(template *models.Template, packagerDir string) error {
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

func (fs *FileService) GetPackagerRoot(library string) (string, error) {
	templateDir, err := fs.GetTemplatesDir()
	if err != nil {
		return "", err
	}
	templatePath := filepath.Join(templateDir, library+".json")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	template, err := utils.ParseTemplate(content)

	if err != nil {
		return "", err
	}

	buildFile := template.Detect["manifest"]
	if buildFile == "" {
		return "", errors.New("no manifest field found in template -> 'detect'")
	}

	rootDir, err := FindProjectRoot(templatePath, buildFile)
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
