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

func (s *SearchService) Search(query string) ([]map[string]interface{}, error) {
	return s.repo.Search("emails", query)
}
