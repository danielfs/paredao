package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danielfs/paredao/backend/entities"
	"github.com/danielfs/paredao/backend/repositories"
	"github.com/gorilla/mux"
)

// GetVotos handles GET /votos
func GetVotos(w http.ResponseWriter, r *http.Request) {
	votos := repositories.GetAllVotos()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votos)
}

// GetVoto handles GET /votos/{participanteId}/{votacaoId}
func GetVoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	participanteId, err := strconv.ParseInt(vars["participanteId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid participanteId format", http.StatusBadRequest)
		return
	}

	votacaoId, err := strconv.ParseInt(vars["votacaoId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid votacaoId format", http.StatusBadRequest)
		return
	}

	voto, exists := repositories.GetVotoByIDs(participanteId, votacaoId)
	if !exists {
		http.Error(w, "Voto not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(voto)
}

// CreateVoto handles POST /votos
func CreateVoto(w http.ResponseWriter, r *http.Request) {
	var votoRequest struct {
		ParticipanteId int64 `json:"participanteId"`
		VotacaoId      int64 `json:"votacaoId"`
	}

	err := json.NewDecoder(r.Body).Decode(&votoRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if votoRequest.ParticipanteId == 0 || votoRequest.VotacaoId == 0 {
		http.Error(w, "ParticipanteId and VotacaoId are required", http.StatusBadRequest)
		return
	}

	// Check if participante exists
	participante, exists := repositories.GetParticipanteByID(votoRequest.ParticipanteId)
	if !exists {
		http.Error(w, "Participante not found", http.StatusNotFound)
		return
	}

	// Check if votacao exists
	votacao, exists := repositories.GetVotacaoByID(votoRequest.VotacaoId)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Create voto
	voto := &entities.Voto{
		Participante: participante,
		Votacao:      votacao,
	}

	// Save voto
	savedVoto := repositories.SaveVoto(voto)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedVoto)
}
