package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"project/internal/repository"
	"project/internal/services"
)

type SearchHandler struct {
	service *services.SearchService
}

func NewSearchHandler(service *services.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Parámetros de búsqueda
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query requerida", http.StatusBadRequest)
		return
	}

	// Parámetros opcionales
	from, _ := strconv.Atoi(r.URL.Query().Get("from"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	// Campos de búsqueda
	fields := r.URL.Query()["fields"]

	// Opciones de búsqueda
	searchOptions := repository.SearchOptions{
		Query:     query,
		Fields:    fields,
		From:      from,
		Size:      size,
		SortField: sortField,
		SortOrder: sortOrder,
	}

	// Realizar búsqueda
	results, err := h.service.Search(query, searchOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
func (h *SearchHandler) ListEmails(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de paginación
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	// Valores por defecto
	page := 0
	size := 100

	// Convertir página si se proporciona
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	// Convertir tamaño si se proporciona
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil {
			size = s
		}
	}

	// Configurar opciones de búsqueda
	options := repository.SearchOptions{
		From: page * size,
		Size: size,
	}

	// Realizar búsqueda
	results, err := h.service.ListEmails(options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listando correos: %v", err), http.StatusInternalServerError)
		return
	}

	// Devolver resultados
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
