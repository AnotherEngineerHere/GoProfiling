package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"project/internal/handlers"
	"project/internal/repository"
	"project/internal/services"
)

func main() {
	port := flag.String("port", "3000", "Puerto para el servidor")
	flag.Parse()

	// Inicializar router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Permite todos los orígenes
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Tiempo máximo de caché para preflight requests
	}))

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
		r.Get("/emails", searchHandler.ListEmails) // Nuevo endpoint para listar correos
	})

	// Iniciar servidor
	log.Printf("Servidor corriendo en puerto %s", *port)
	http.ListenAndServe(":"+*port, r)
}
