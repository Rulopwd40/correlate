package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateTaskLog(taskName, workdir string) (io.Writer, error) {
	if workdir == "" {
		return nil, fmt.Errorf("workdir no puede estar vacío")
	}

	info, err := os.Stat(workdir)
	if err != nil {
		return nil, fmt.Errorf("workdir no existe: %s", workdir)
	}
	if !info.IsDir() {
		workdir = filepath.Dir(workdir)
	}
	finalPath := filepath.Join(workdir, ".correlate/logs")
	err = os.MkdirAll(finalPath, info.Mode())
	if err != nil {
		return nil, fmt.Errorf("workdir no existe: %s", workdir)
	}
	logFilePath := filepath.Join(finalPath, taskName+".log")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear log file: %v", err)
	}

	return logFile, nil
}
