package core

import (
	"errors"
	"github.com/Rulopwd40/correlate/internal/models"
	"log"
)

const CORRELATE_PATH = ".correlate"

type Orchestrator struct {
	fs  iFileService
	cfg iConfigService
	ts  iTemplateService
	rs  iReferenceService
}

func NewOrchestrator(fs iFileService, cfg iConfigService, ts iTemplateService, rs iReferenceService) *Orchestrator {
	return &Orchestrator{
		fs:  fs,
		cfg: cfg,
		ts:  ts,
		rs:  rs,
	}
}

/*
Library: indicates template name for a library
name: indicates project identifier
*/
func (orch *Orchestrator) Init(library string, name string) error {

	cfg, err := orch.cfg.GenerateConfig(library, name)
	if err != nil {
		log.Println("Error generating config:", err.Error())
		return err
	}
	_, err = orch.ts.GenerateProjectTemplate(library, name, cfg.PackageDirectory)
	if err != nil {
		log.Println("Error generating project template:", err.Error())
		return err
	}

	err = orch.rs.GenerateReferencesFile()
	if err != nil {
		log.Println("Error generating references file:", err.Error())
		return err
	}

	return nil
}

func (orch *Orchestrator) Link(identifier string, projectRoot string) error {
	template, err := orch.ts.GetProjectTemplate()
	if err != nil {
		log.Println("Error getting project template:", err.Error())
		return err
	}

	manifest := template.Detect["manifest"]
	if manifest == "" {
		log.Println("Error detecting manifest")
		return errors.New("Error detecting manifest")
	}

	log.Printf("Looking for all occurrences of %s in all '%s' files of the project, this might take a while.", identifier, manifest)
	paths, err := orch.fs.FindAllFilesThatContains(projectRoot, manifest, identifier)
	if err != nil {
		log.Println("Error finding files:", err.Error())
		return err
	}
	reference := models.Reference{
		Identifier:          identifier,
		ManifestDirectories: []string{},
	}

	for _, path := range paths {
		reference.ManifestDirectories = append(reference.ManifestDirectories, path)
	}

	references := models.References{
		append([]models.Reference{}, reference),
	}

	err = orch.rs.SaveReferences(references)
	if err != nil {
		log.Println("Error saving references:", err.Error())
		return err
	}
	return nil
}
