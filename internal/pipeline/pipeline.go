package pipeline

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Rulopwd40/correlate/internal/logger"
)

type Pipeline struct {
	Tasks      []Task
	WorkingDir string
	EventSink  chan Event
}

func (p *Pipeline) Run(ctx context.Context) {
	for _, task := range p.Tasks {
		p.emit(Event{
			Type:     EventTaskStart,
			TaskName: task.Name,
			Message:  "starting",
		})

		if err := p.runTask(ctx, task); err != nil {
			p.emit(Event{
				Type:     EventError,
				TaskName: task.Name,
				Message:  err.Error(),
				Err:      err,
			})
			return
		}

		p.emit(Event{
			Type:     EventTaskFinish,
			TaskName: task.Name,
			Message:  "completed",
		})
	}

	// Avisar que el pipeline completo terminó
	p.emit(Event{Type: EventPipelineDone})
}

func (p *Pipeline) emit(ev Event) {
	if p.EventSink != nil {
		p.EventSink <- ev
	}
}

func (p *Pipeline) runTask(ctx context.Context, t Task) error {
	var cmd *exec.Cmd
	workdir := t.Workdir
	if workdir == "" {
		workdir = p.WorkingDir
	}
	info, err := os.Stat(workdir)
	if err != nil {
		return fmt.Errorf("workdir does not exist: %s", workdir)
	}
	if !info.IsDir() {
		workdir = filepath.Dir(workdir)
	}

	// Detectar si estamos en Windows
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", t.Cmd)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", t.Cmd)
	}
	cmd.Dir = workdir

	logWriter, err := logger.CreateTaskLog(t.Name, workdir)
	if err != nil {
		return err
	}
	defer func() {
		if c, ok := logWriter.(io.Closer); ok {
			c.Close()
		}
	}()

	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
func (p *Pipeline) stream(prefix string, taskName string, pipe interface{}) {
	var r *bufio.Reader

	switch v := pipe.(type) {
	case *bufio.Reader:
		r = v
	default:
		r = bufio.NewReader(v.(interface{ Read([]byte) (int, error) }))
	}

	for {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			p.emit(Event{
				Type:     EventTaskProgress,
				TaskName: taskName,
				Message:  line,
			})
		}
		if err != nil {
			return
		}
	}
}
