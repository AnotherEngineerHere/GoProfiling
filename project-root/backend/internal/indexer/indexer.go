package indexer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"project/internal/repository"
)

type EmailIndexer struct {
	repo *repository.ZincRepository
}

func NewEmailIndexer(repo *repository.ZincRepository) *EmailIndexer {
	return &EmailIndexer{repo: repo}
}

func (e *EmailIndexer) IndexEmailsFromPath(path string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(filePath, ".txt") {
			return e.indexSingleEmail(filePath)
		}

		return nil
	})
}

func (e *EmailIndexer) indexSingleEmail(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error leyendo archivo %s: %v", filePath, err)
		return nil
	}

	email := map[string]interface{}{
		"path":    filePath,
		"content": string(content),
	}

	err = e.repo.Index("emails", email)
	if err != nil {
		log.Printf("Error indexando %s: %v", filePath, err)
		return nil
	}

	fmt.Printf("Indexado: %s\n", filePath)
	return nil
}
