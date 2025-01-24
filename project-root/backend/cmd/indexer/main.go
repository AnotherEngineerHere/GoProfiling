package main

import (
	"fmt"
	"log"
	"os"

	"project/internal/indexer"
	"project/internal/repository"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Por favor proporciona la ruta del dataset")
	}

	datasetPath := os.Args[1]

	// Configurar repositorio de ZincSearch
	zincRepo, err := repository.NewZincRepository()
	if err != nil {
		log.Fatalf("Error inicializando ZincSearch: %v", err)
	}

	// Crear indexador
	idx := indexer.NewEmailIndexer(zincRepo)

	// Indexar correos
	err = idx.IndexEmailsFromPath(datasetPath)
	if err != nil {
		log.Fatalf("Error indexando correos: %v", err)
	}

	fmt.Println("Indexación completada exitosamente")
}
