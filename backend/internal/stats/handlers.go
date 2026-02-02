package stats

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ngthecoder/go_web_api/internal/errors"
)

type StatsHandler struct {
	service *StatsService
}

func NewStatsHandler(service *StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

func (h *StatsHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	categoryCounts, err := h.service.GetCategoryCounts()
	if err != nil {
		log.Printf("Error getting category counts: %v", err)
		errors.WriteHTTPError(w, errors.NewInternalServerError("Failed to get categories", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categoryCounts); err != nil {
		log.Printf("Error encoding category counts response: %v", err)
		errors.WriteHTTPError(w, errors.NewInternalServerError("Failed to encode response", err))
		return
	}
}

func (h *StatsHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	stats, err := h.service.GetStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		errors.WriteHTTPError(w, errors.NewInternalServerError("Failed to get stats", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding stats response: %v", err)
		errors.WriteHTTPError(w, errors.NewInternalServerError("Failed to encode response", err))
		return
	}
}
