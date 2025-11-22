package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Rulopwd40/correlate/internal/models"
	"github.com/Rulopwd40/correlate/internal/utils"
)

type iConfigService interface {
	GenerateConfig(library, identifier string) (models.Config, error)
	SaveConfig(config models.Config) error
	GetConfig() (models.Config, error)
}
type ConfigService struct {
	fs  iFileService
	ts  iTemplateService
	env *Environment
}

func NewConfigService(fs iFileService, ts iTemplateService) *ConfigService {
	env := DefaultEnvironment()
	return &ConfigService{
		fs:  fs,
		ts:  ts,
		env: env,
	}
}

// Generate Config.json
func (cf *ConfigService) GenerateConfig(library, identifier string) (models.Config, error) {
	templateExists := cf.ts.TemplateExists(library)
	if !templateExists {
		return models.Config{}, errors.New("template does not exist")
	}
	config := models.Config{
		TemplateName: library,
		Variables:    map[string]string{"identifier": identifier},
	}

	var err error
	log.Println("Getting Packager")
	config.PackageDirectory, err = cf.ts.GetPackagerRoot(library)
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
	dir := cf.env.GetCorrelatePath()

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
func (cf *ConfigService) GetConfig() (models.Config, error) {
	filePath := filepath.Join(cf.env.GetCorrelatePath(), "config.json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return models.Config{}, fmt.Errorf("config.json not found at %s", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return models.Config{}, err
	}

	config, err := utils.ParseConfig(data)
	if err != nil {
		return models.Config{}, err
	}

	return config, nil
}
