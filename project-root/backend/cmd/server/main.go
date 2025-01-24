package main

import (
	"flag"
	"log"
	"net/http"
	"project/internal/handlers"
	"project/internal/repository"
	"project/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port := flag.String("port", "3000", "Puerto para el servidor")
	flag.Parse()

	// Inicializar router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Inicializar repositorio
	zincRepo, err := repository.NewZincRepository()
	if err != nil {
		log.Fatalf("Error inicializando ZincSearch: %v", err)
	}

	// Inicializar servicio de búsqueda
	searchService := services.NewSearchService(zincRepo)

	// Inicializar handlers
	searchHandler := handlers.NewSearchHandler(searchService)

	// Definir rutas
	r.Route("/api", func(r chi.Router) {
		r.Get("/search", searchHandler.Search)
	})

	// Iniciar servidor
	log.Printf("Servidor corriendo en puerto %s", *port)
	http.ListenAndServe(":"+*port, r)
}
