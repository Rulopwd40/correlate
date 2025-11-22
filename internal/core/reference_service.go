package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Rulopwd40/correlate/internal/models"
	"github.com/Rulopwd40/correlate/internal/utils"
)

type iReferenceService interface {
	GenerateReferencesFile() error
	SaveReferences(ref models.References) error
	GetReferences() (models.References, error)
	GetReference(identifier string) (models.Reference, error)
}

type ReferenceService struct {
	fs  iFileService
	env *Environment
}

func NewReferenceService(fs iFileService) *ReferenceService {
	env := DefaultEnvironment()
	return &ReferenceService{fs, env}
}

func (rs *ReferenceService) GenerateReferencesFile() error {
	reference := models.References{
		References: []models.Reference{},
	}

	err := rs.SaveReferences(reference)
	return err
}

func (rs *ReferenceService) SaveReferences(ref models.References) error {
	log.Println("Writing references file")
	dir := rs.env.GetCorrelatePath()
	filePath := filepath.Join(dir, "references.json")

	// Asegurar directorio
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fileReferences, err := rs.GetReferences()
	if err != nil {
		return err
	}

	for _, incomingRef := range ref.References {

		existingIndex := -1
		for i, ex := range fileReferences.References {
			if ex.Identifier == incomingRef.Identifier {
				existingIndex = i
				break
			}
		}

		if existingIndex != -1 {
			merged := mergeUniquePaths(
				fileReferences.References[existingIndex].ManifestDirectories,
				incomingRef.ManifestDirectories,
			)
			fileReferences.References[existingIndex].ManifestDirectories = merged

		} else {
			// No existe ese identifier → agregarlo entero
			fileReferences.References = append(
				fileReferences.References,
				incomingRef,
			)
		}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(fileReferences); err != nil {
		return err
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return err
	}

	fmt.Println("references.json generated successfully.")
	return nil
}

func (rs *ReferenceService) GetReferences() (models.References, error) {
	dir := rs.env.GetCorrelatePath()
	filePath := filepath.Join(dir, "references.json")
	var fileReferences models.References

	data, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return models.References{}, err
		}
		fileReferences = models.References{
			References: []models.Reference{},
		}
	} else {
		fileReferences, err = utils.ParseReferences(data)
		if err != nil {
			return models.References{}, err
		}
	}
	return fileReferences, nil
}

func (rs *ReferenceService) GetReference(identifier string) (models.Reference, error) {
	references, err := rs.GetReferences()
	if err != nil {
		return models.Reference{}, err
	}
	for _, ref := range references.References {
		if ref.Identifier == identifier {
			return ref, nil
		}
	}
	return models.Reference{}, errors.New("no reference found with identifier: " + identifier)
}

func mergeUniquePaths(existing, incoming []string) []string {
	seen := make(map[string]bool)

	for _, p := range existing {
		seen[p] = true
	}

	for _, p := range incoming {
		if !seen[p] {
			existing = append(existing, p)
			seen[p] = true
		}
	}

	return existing
}
