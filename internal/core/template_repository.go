package core

import (
	"fmt"
	"io"
	"net/http"
)

type TemplateRepository interface {
	FetchTemplate(name string) ([]byte, error)
}

type RestTemplateRepository struct {
	BaseURL string
}

func NewRestTemplateRepository(baseURL string) *RestTemplateRepository {
	return &RestTemplateRepository{BaseURL: baseURL}
}

func (r *RestTemplateRepository) FetchTemplate(name string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s.json", r.BaseURL, name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
