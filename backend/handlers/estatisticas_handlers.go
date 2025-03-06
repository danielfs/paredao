package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/danielfs/paredao/backend/entities"
	"github.com/danielfs/paredao/backend/repositories"
)

func getVotacaoData(
	w http.ResponseWriter,
	r *http.Request,
	cacheKeyFormat string,
	fetchData func(int64) (interface{}, error),
	errorMsg string,
) {
	vars := mux.Vars(r)
	ctx := r.Context()

	votacaoID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid votacaoID format", http.StatusBadRequest)
		return
	}

	_, exists := repositories.GetVotacaoByID(votacaoID)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	cacheKey := fmt.Sprintf(cacheKeyFormat, votacaoID)

	var data interface{}
	found, err := repositories.GetFromCache(ctx, cacheKey, &data)
	if err != nil {
		// Registra o erro mas continua com a consulta ao banco de dados
		fmt.Printf("Cache error: %v\n", err)
	}

	if !found {
		// Cache não encontrado, busca no banco de dados
		data, err = fetchData(votacaoID)
		if err != nil {
			http.Error(w, errorMsg, http.StatusInternalServerError)
			return
		}

		// Armazena no cache para requisições futuras
		if err := repositories.SetCache(ctx, cacheKey, data); err != nil {
			// Registra o erro mas continua mesmo se o cache falhar
			fmt.Printf("Cache set error: %v\n", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Printf("JSON encode error: %v\n", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetVotacaoTotal(w http.ResponseWriter, r *http.Request) {
	getVotacaoData(
		w,
		r,
		repositories.TotalCacheKey,
		func(votacaoID int64) (interface{}, error) {
			total, err := repositories.GetTotalVotesForVotacao(votacaoID)
			if err != nil {
				return nil, err
			}
			return entities.VotacaoTotalResponse{
				VotacaoID: votacaoID,
				Total:     total,
			}, nil
		},
		"Error getting total votes",
	)
}

func GetVotacaoTotalByParticipante(w http.ResponseWriter, r *http.Request) {
	getVotacaoData(
		w,
		r,
		repositories.ParticipantCacheKey,
		func(votacaoID int64) (interface{}, error) {
			return repositories.GetTotalVotesByParticipante(votacaoID)
		},
		"Error getting total votes by participante",
	)
}

func GetVotacaoTotalByHour(w http.ResponseWriter, r *http.Request) {
	getVotacaoData(
		w,
		r,
		repositories.HourlyCacheKey,
		func(votacaoID int64) (interface{}, error) {
			return repositories.GetTotalVotesByHour(votacaoID)
		},
		"Error getting total votes by hour",
	)
}
