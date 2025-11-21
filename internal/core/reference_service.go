package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Rulopwd40/correlate/internal/models"
	"github.com/Rulopwd40/correlate/internal/utils"
	"log"
	"os"
	"path/filepath"
)

type iReferenceService interface {
	GenerateReferencesFile() error
	SaveReferences(ref models.References) error
}

type ReferenceService struct {
	fs iFileService
}

func NewReferenceService(fs iFileService) *ReferenceService {
	return &ReferenceService{fs}
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
	dir := CORRELATE_PATH
	filePath := filepath.Join(dir, "references.json")

	// Asegurar directorio
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var fileReferences models.References

	data, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		fileReferences = models.References{
			References: []models.Reference{},
		}
	} else {
		fileReferences, err = utils.ParseReferences(data)
		if err != nil {
			return err
		}
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
