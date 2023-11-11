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
	r.ParseForm()
	q := r.Form["q"]
	results, err := h.index.Search(q...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := handlers_models.SearchResponse{Results: results}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
