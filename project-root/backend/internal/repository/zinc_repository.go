package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ZincRepository struct {
	baseURL string
	client  *http.Client
}

func NewZincRepository() (*ZincRepository, error) {
	return &ZincRepository{
		baseURL: "http://localhost:4080", // URL por defecto de ZincSearch
		client:  &http.Client{},
	}, nil
}

func (z *ZincRepository) Index(index string, document interface{}) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/%s/_doc", z.baseURL, index),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth("admin", "admin") // Credenciales por defecto
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error indexando: %s", string(body))
	}

	return nil
}

func (z *ZincRepository) Search(index, query string) ([]map[string]interface{}, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"term": query,
		},
	}

	jsonData, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/%s/_search", z.baseURL, index),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("admin", "admin")
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no se encontraron resultados")
	}

	return hits["hits"].([]map[string]interface{}), nil
}
