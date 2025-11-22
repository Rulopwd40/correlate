package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type iFileService interface {
	GetJSON(path string, name string) ([]byte, error)
	GetFile(path string, name string) ([]byte, error)
	FindAllFiles(root string, name string) ([]string, error)
	FindAllFilesThatContains(root string, filename string, identifier string) ([]string, error)
	GetRegexOcurrenceAfter(path string, regex string, identifier string) (string, error)
	ReplaceRegexOccurrenceAfter(path string, regex string, identifier string, replacement string) (string, error)
}

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (fs *FileService) GetJSON(path string, name string) ([]byte, error) {
	fullPath, err := fs.resolvePath(path, name+".json")
	if err != nil {
		return nil, err
	}

	return os.ReadFile(fullPath)
}
func (fs *FileService) GetFile(path string, name string) ([]byte, error) {
	fullPath, err := fs.resolvePath(path, name)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(fullPath)
}

func (fs *FileService) FindAllFiles(root string, name string) ([]string, error) {
	var results []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // error al acceder a un archivo/directorio
		}

		if !info.IsDir() && info.Name() == name {
			results = append(results, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (fs *FileService) FindAllFilesThatContains(root string, filename string, identifier string) ([]string, error) {
	var results []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == filename {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if bytes.Contains(data, []byte(identifier)) {
				results = append(results, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (fs *FileService) GetRegexOcurrenceAfter(path string, regex string, identifier string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	content := string(data)

	idx := strings.Index(content, identifier)
	if idx == -1 {
		return "", fmt.Errorf("identifier %q not found", identifier)
	}

	sub := content[idx+len(identifier):]

	re, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}

	match := re.FindString(sub)
	if match == "" {
		return "", fmt.Errorf("no regex match found after identifier")
	}

	return match, nil
}
func (fs *FileService) ReplaceRegexOccurrenceAfter(path string, regex string, identifier string, replacement string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := string(data)

	idx := strings.Index(content, identifier)
	if idx == -1 {
		return "", fmt.Errorf("identifier %q not found", identifier)
	}
	before := content[:idx+len(identifier)]
	after := content[idx+len(identifier):]

	re, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}

	loc := re.FindStringIndex(after)
	if loc == nil {
		return "", fmt.Errorf("regex %q not found after identifier", regex)
	}

	original := after[loc[0]:loc[1]]

	updatedAfter := after[:loc[0]] + replacement + after[loc[1]:]

	finalContent := before + updatedAfter

	err = os.WriteFile(path, []byte(finalContent), 0644)
	if err != nil {
		return "", err
	}

	return original, nil
}

func (fs *FileService) resolvePath(path, name string) (string, error) {
	base, err := fs.baseDirFrom(path)
	if err != nil {
		return "", err
	}

	full := filepath.Join(base, name)
	return full, nil
}

func (fs *FileService) baseDirFrom(path string) (string, error) {
	if path == "" {
		return os.Getwd()
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, path), nil
}
