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

func (z *ZincRepository) CreateIndex(index string) error {
	indexConfig := map[string]interface{}{
		"name": index,
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"path":      map[string]string{"type": "keyword"},
				"content":   map[string]string{"type": "text"},
				"sender":    map[string]string{"type": "keyword"},
				"recipient": map[string]string{"type": "keyword"},
				"subject":   map[string]string{"type": "text"},
				"date":      map[string]string{"type": "date"},
				"folder":    map[string]string{"type": "keyword"},
			},
		},
	}

	jsonData, err := json.Marshal(indexConfig)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/index", z.baseURL),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth("admin", "admin")
	req.Header.Set("Content-Type", "application/json")

	resp, err := z.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error creando índice: %s", string(body))
	}

	return nil
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

	// Leer el cuerpo de la respuesta para obtener más información
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error indexando: %s (código de estado: %d)", string(body), resp.StatusCode)
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

type SearchOptions struct {
	Query     string
	Fields    []string
	From      int
	Size      int
	SortField string
	SortOrder string
}

func (z *ZincRepository) ListAll(index string, opts SearchOptions) (map[string]interface{}, error) {
	// Consulta para recuperar todos los documentos
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"from": opts.From,
		"size": opts.Size,
		"sort": []map[string]interface{}{
			{
				"date": map[string]interface{}{
					"order": "desc", // Ordenar por fecha más reciente primero
				},
			},
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

	return result, nil
}

func (z *ZincRepository) AdvancedSearch(index string, opts SearchOptions) (map[string]interface{}, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  opts.Query,
				"fields": opts.Fields,
			},
		},
		"from": opts.From,
		"size": opts.Size,
	}

	// Añadir ordenamiento si se especifica
	if opts.SortField != "" {
		searchQuery["sort"] = []map[string]interface{}{
			{
				opts.SortField: map[string]interface{}{
					"order": opts.SortOrder,
				},
			},
		}
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

	return result, nil
}
