package main

import (
	"github.com/Rulopwd40/correlate/internal/commands"
	"github.com/Rulopwd40/correlate/internal/core"
)

func main() {
	initialize()
	commands.Execute()

}

func initialize() {
	fs := core.NewFileService()
	cs := core.NewConfigService(fs)
	ts := core.NewTemplateService(fs)
	rs := core.NewReferenceService(fs)

	orch := core.NewOrchestrator(fs, cs, ts, rs)

	// Registrar en el contexto
	core.Register(fs)
	core.Register(cs)
	core.Register(ts)
	core.Register(rs)
	core.Register(orch)
}
