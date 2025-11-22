package core

import (
	"os"
	"path/filepath"
)

// Environment holds all configuration constants for the application
type Environment struct {
	TemplatesPath         string
	ProjectTemplatesPath  string
	CorrelatePath         string
	TemplateRepositoryURL string
}

// DefaultEnvironment returns the default environment configuration
func DefaultEnvironment() *Environment {
	return &Environment{
		TemplatesPath:         getEnvOrDefault("CORRELATE_TEMPLATES_PATH", "templates"),
		ProjectTemplatesPath:  getEnvOrDefault("CORRELATE_PROJECT_TEMPLATES_PATH", ".correlate"),
		CorrelatePath:         getEnvOrDefault("CORRELATE_PATH", ".correlate"),
		TemplateRepositoryURL: getEnvOrDefault("CORRELATE_TEMPLATE_REPOSITORY_URL", "https://raw.githubusercontent.com/Rulopwd40/correlate/develop/templates"),
	}
}

// getEnvOrDefault retrieves an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetTemplatesPath returns the absolute path to the templates directory
func (e *Environment) GetTemplatesPath() string {
	if filepath.IsAbs(e.TemplatesPath) {
		return e.TemplatesPath
	}
	return e.TemplatesPath
}

// GetProjectTemplatesPath returns the path to project-specific templates
func (e *Environment) GetProjectTemplatesPath() string {
	return e.ProjectTemplatesPath
}

// GetCorrelatePath returns the path to the correlate working directory
func (e *Environment) GetCorrelatePath() string {
	return e.CorrelatePath
}

// GetTemplateRepositoryURL returns the URL to the template repository
func (e *Environment) GetTemplateRepositoryURL() string {
	return e.TemplateRepositoryURL
}
