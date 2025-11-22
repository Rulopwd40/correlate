package pipeline

import (
	"bufio"
	"io"
	"os/exec"
	"runtime"
)

type Executor struct {
	Events chan Event
}

func NewExecutor() *Executor {
	return &Executor{
		Events: make(chan Event, 10),
	}
}

func (e *Executor) RunPipeline(p Pipeline) {
	go func() {
		for _, task := range p.Tasks {
			e.runTask(task)
		}

		e.Events <- Event{Type: EventPipelineDone}
		close(e.Events)
	}()
}

func (e *Executor) runTask(t Task) {
	e.Events <- Event{
		Type:     EventTaskStart,
		TaskName: t.Name,
		Message:  "Started",
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", t.Cmd)
	} else {
		cmd = exec.Command("sh", "-c", t.Cmd) // <-- ESTE es el correcto
	}

	cmd.Dir = t.Workdir

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	_ = cmd.Start()

	go e.stream(stdout, t)
	go e.stream(stderr, t)

	_ = cmd.Wait()

	e.Events <- Event{
		Type:     EventTaskFinish,
		TaskName: t.Name,
		Message:  "Finished",
	}
}

func (e *Executor) stream(pipe io.ReadCloser, t Task) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		e.Events <- Event{
			Type:     EventTaskProgress,
			TaskName: t.Name,
			Message:  scanner.Text(),
		}
	}
}
