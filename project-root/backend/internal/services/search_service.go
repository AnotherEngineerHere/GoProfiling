package services

import (
	"project/internal/repository"
)

type SearchService struct {
	repo *repository.ZincRepository
}

func NewSearchService(repo *repository.ZincRepository) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) Search(query string, options repository.SearchOptions) (map[string]interface{}, error) {
	// Configurar opciones por defecto si no se proporcionan
	if options.Fields == nil {
		options.Fields = []string{"content", "sender", "recipient", "subject", "folder"}
	}
	if options.From == 0 {
		options.From = 0
	}
	if options.Size == 0 {
		options.Size = 10
	}
	options.Query = query

	return s.repo.AdvancedSearch("emails", options)
}
func (s *SearchService) ListEmails(options repository.SearchOptions) (map[string]interface{}, error) {
	// Configurar opciones por defecto
	if options.From == 0 {
		options.From = 0
	}
	if options.Size == 0 {
		options.Size = 100 // Número predeterminado de correos a mostrar
	}

	return s.repo.ListAll("emails", options)
}
