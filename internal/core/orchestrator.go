package core

import (
	stdcontext "context"
	"errors"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Rulopwd40/correlate/internal/models"
	"github.com/Rulopwd40/correlate/internal/pipeline"
	"github.com/Rulopwd40/correlate/internal/utils"
)

type Orchestrator struct {
	fs  iFileService
	cfg iConfigService
	ts  iTemplateService
	rs  iReferenceService
	env *Environment

	eventStream chan pipeline.Event
}

func NewOrchestrator(fs iFileService, cfg iConfigService, ts iTemplateService, rs iReferenceService) *Orchestrator {
	env := DefaultEnvironment()
	return &Orchestrator{
		fs:  fs,
		cfg: cfg,
		ts:  ts,
		rs:  rs,
		env: env,
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
		log.Println("Error: manifest file not specified in template")
		return errors.New("manifest file not specified in template")
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
		References: append([]models.Reference{}, reference),
	}

	err = orch.rs.SaveReferences(references)
	if err != nil {
		log.Println("Error saving references:", err.Error())
		return err
	}
	return nil
}

func (orch *Orchestrator) Replace(identifier string, version string) error {
	template, err := orch.ts.GetProjectTemplate()
	if err != nil {
		log.Println("Error getting project template:", err.Error())
		return err
	}

	cfg, err := orch.cfg.GetConfig()
	if err != nil {
		log.Println("Error getting config:", err.Error())
		return err
	}

	manifest := template.Detect["manifest"]
	versionPattern := template.Variables["version"]
	versionRegex := strings.Replace(versionPattern, "{{version}}", `([0-9A-Za-z\.\-\_]+)`, -1)
	replaceVersion := version

	if replaceVersion == "" {
		reference, err := orch.rs.GetReference(identifier)
		if err != nil {
			log.Println("Error getting reference:", err.Error())
			return err
		}
		log.Printf("Not version specified, searching in target manifest %s %s %s", identifier, manifest, versionRegex)
		var versions []string
		for _, path := range reference.ManifestDirectories {
			projectVersion, err := orch.fs.GetRegexOcurrenceAfter(path, versionRegex, identifier)
			if err != nil {
				log.Printf("Error finding project version for %s: %s", identifier, path)
				return err
			}
			versions = append(versions, projectVersion)
		}

		if len(versions) == 0 {
			return errors.New("Error finding project version. Look references.json for " + identifier)
		}
		if !utils.AllEqual(versions) {
			return errors.New("error, all project versions must be equals")
		}

		replaceVersion = versions[0]
	} else {
		replaceVersion = strings.Replace(versionRegex, "{{version}}", version, -1)
	}
	log.Printf("Version obtained: %s", replaceVersion)
	manifestPath := filepath.Join(cfg.PackageDirectory, manifest)
	log.Printf("Replacing version in %s", manifestPath)
	_, err = orch.fs.ReplaceRegexOccurrenceAfter(manifestPath, versionRegex, identifier, replaceVersion)
	if err != nil {
		return err
	}
	return nil
}

func (orch *Orchestrator) Update(identifier string) error {

	template, err := orch.ts.GetProjectTemplate()
	if err != nil {
		return err
	}

	cfg, err := orch.cfg.GetConfig()
	if err != nil {
		return err
	}

	orch.eventStream = make(chan pipeline.Event)

	var references []models.Reference
	if identifier == "" {
		ref, err := orch.rs.GetReferences()
		if err != nil {
			return err
		}
		references = ref.References
	} else {
		ref, err := orch.rs.GetReference(identifier)
		if err != nil {
			return err
		}
		references = []models.Reference{ref}
	}

	if references == nil || len(references) == 0 {
		return errors.New("Error finding references")
	}

	go func() {
		wg := sync.WaitGroup{}

		for _, ref := range references {
			log.Printf("Generating pipeline for %s", ref.Identifier)
			for _, dir := range ref.ManifestDirectories {

				ref := ref
				dir := dir

				wg.Add(1)

				go func() {
					defer wg.Done()

					// Build context with all available paths
					context := map[string]string{
						"sourceDir":         cfg.PackageDirectory, // Where correlate init was run (the library being developed)
						"targetDir":         dir,                  // Where the consumer/dependent project is
						"projectIdentifier": ref.Identifier,       // The identifier from references.json
					}

					tasks, err := utils.MakeTasks(template.Steps, ref, context)
					if err != nil {
						orch.eventStream <- pipeline.Event{
							Type:    pipeline.EventError,
							Err:     err,
							Message: err.Error(),
						}
						return
					}

					p := &pipeline.Pipeline{
						Tasks:      tasks,
						WorkingDir: dir,
						EventSink:  orch.eventStream,
					}

					p.Run(stdcontext.Background())
				}()
			}
		}

		wg.Wait()
		close(orch.eventStream)
	}()

	return nil
}

func (orch *Orchestrator) Events() <-chan pipeline.Event {
	return orch.eventStream
}
