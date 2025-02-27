package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danielfs/paredao/backend/entities"
	"github.com/danielfs/paredao/backend/repositories"
	"github.com/gorilla/mux"
)

// GetVotacaoTotal handles GET /estatisticas/votacoes/{id}/total
func GetVotacaoTotal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()

	votacaoId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid votacaoId format", http.StatusBadRequest)
		return
	}

	// Check if votacao exists
	_, exists := repositories.GetVotacaoByID(votacaoId)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Create cache key
	cacheKey := fmt.Sprintf(repositories.TotalCacheKey, votacaoId)

	// Try to get from cache first
	var response entities.VotacaoTotalResponse
	found, err := repositories.GetFromCache(ctx, cacheKey, &response)
	if err != nil {
		// Continue with database query on cache error
	}

	if !found {
		// Cache miss, get from database
		total, err := repositories.GetTotalVotesForVotacao(votacaoId)
		if err != nil {
			http.Error(w, "Error getting total votes", http.StatusInternalServerError)
			return
		}

		response = entities.VotacaoTotalResponse{
			VotacaoId: votacaoId,
			Total:     total,
		}

		// Store in cache for future requests
		if err := repositories.SetCache(ctx, cacheKey, response); err != nil {
			// Continue even if caching fails
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetVotacaoTotalByParticipante handles GET /estatisticas/votacoes/{id}/participantes
func GetVotacaoTotalByParticipante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()

	votacaoId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid votacaoId format", http.StatusBadRequest)
		return
	}

	// Check if votacao exists
	_, exists := repositories.GetVotacaoByID(votacaoId)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Create cache key
	cacheKey := fmt.Sprintf(repositories.ParticipantCacheKey, votacaoId)

	// Try to get from cache first
	var totals []entities.ParticipanteTotalResponse
	found, err := repositories.GetFromCache(ctx, cacheKey, &totals)
	if err != nil {
		// Continue with database query on cache error
	}

	if !found {
		// Cache miss, get from database
		totals, err = repositories.GetTotalVotesByParticipante(votacaoId)
		if err != nil {
			http.Error(w, "Error getting total votes by participante", http.StatusInternalServerError)
			return
		}

		// Store in cache for future requests
		if err := repositories.SetCache(ctx, cacheKey, totals); err != nil {
			// Continue even if caching fails
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totals)
}

// GetVotacaoTotalByHour handles GET /estatisticas/votacoes/{id}/hourly
func GetVotacaoTotalByHour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()

	votacaoId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid votacaoId format", http.StatusBadRequest)
		return
	}

	// Check if votacao exists
	_, exists := repositories.GetVotacaoByID(votacaoId)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Create cache key
	cacheKey := fmt.Sprintf(repositories.HourlyCacheKey, votacaoId)

	// Try to get from cache first
	var totals []entities.HourlyTotalResponse
	found, err := repositories.GetFromCache(ctx, cacheKey, &totals)
	if err != nil {
		// Continue with database query on cache error
	}

	if !found {
		// Cache miss, get from database
		totals, err = repositories.GetTotalVotesByHour(votacaoId)
		if err != nil {
			http.Error(w, "Error getting total votes by hour", http.StatusInternalServerError)
			return
		}

		// Store in cache for future requests
		if err := repositories.SetCache(ctx, cacheKey, totals); err != nil {
			// Continue even if caching fails
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totals)
}
