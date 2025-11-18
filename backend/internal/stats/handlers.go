package stats

import (
	"encoding/json"
	"log"
	"net/http"
)

type StatsHandler struct {
	service *StatsService
}

func NewStatsHandler(service *StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

func (h *StatsHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categoryCounts, err := h.service.GetCategoryCounts()
	if err != nil {
		log.Printf("Error getting category counts: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categoryCounts); err != nil {
		log.Printf("Error encoding category counts response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *StatsHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.service.GetStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding stats response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
