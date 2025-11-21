package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rulopwd40/correlate/internal/models"
	"log"
	"os"
	"path/filepath"
)

type iConfigService interface {
	GenerateConfig(library, identifier string) (models.Config, error)
	SaveConfig(config models.Config) error
}
type ConfigService struct {
	fs iFileService
}

func NewConfigService(fs iFileService) *ConfigService {
	return &ConfigService{
		fs: fs,
	}
}

// Generate Config.json
func (cf *ConfigService) GenerateConfig(library, identifier string) (models.Config, error) {
	templateExists := cf.fs.TemplateExists(library)
	if !templateExists {
		return models.Config{}, errors.New("template does not exist")
	}
	config := models.Config{
		TemplateName: library,
		Variables:    map[string]string{"identifier": identifier},
	}

	var err error
	log.Println("Getting Packager")
	config.PackageDirectory, err = cf.fs.GetPackagerRoot(library)
	if err != nil {
		return models.Config{}, err
	}

	err = cf.SaveConfig(config)
	if err != nil {
		return models.Config{}, err
	}
	return config, nil
}

func (cf *ConfigService) SaveConfig(config models.Config) error {
	log.Println("Writing config file")
	dir := CORRELATE_PATH

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dir, "config.json")
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	fmt.Println("config.json generated successfully.")
	return nil
}
