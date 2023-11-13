package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/getnexar/golang-programming-task/doc-index/pkg/index"
	handlers_models "github.com/getnexar/golang-programming-task/doc-index/pkg/server/handlers/models"
	"go.uber.org/zap"
)

type Handlers struct {
	logger *zap.SugaredLogger
	index  index.IndexInterface
}

type JSONSearchQuery struct {
	Keywords []string `json:"keywords"`
}

func NewHandlers(logger *zap.SugaredLogger, index index.IndexInterface) *Handlers {
	return &Handlers{
		index:  index,
		logger: logger,
	}
}

func (h *Handlers) Healthcheck(w http.ResponseWriter, r *http.Request) {
	response := handlers_models.HealthcheckResponse{Status: "OK"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	q := h.getSearchQuery(r)

	if q == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	results := h.index.Search(q...)

	response := handlers_models.SearchResponse{Results: results}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	q := h.getSearchQuery(r)

	if q == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deletedDocumentsCount := h.index.Delete(q...)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(deletedDocumentsCount)
}

func (h *Handlers) getSearchQuery(r *http.Request) []string {
	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		err := r.ParseForm()
		if err != nil {
			return nil
		}

		return r.Form["q"]
	}

	if r.Method == http.MethodPost {
		if r.Header.Get("Content-Type") == "application/json" {
			var searchQuery JSONSearchQuery
			err := json.NewDecoder(r.Body).Decode(&searchQuery)

			if err != nil {
				return nil
			}

			return searchQuery.Keywords
		} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			err := r.ParseForm()

			if err != nil {
				return nil
			}

			return r.Form["q"]
		} else {
			return nil
		}
	}

	return nil
}
