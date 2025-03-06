package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/danielfs/paredao/backend/entities"
	"github.com/danielfs/paredao/backend/repositories"
)

func GetVotos(w http.ResponseWriter, r *http.Request) {
	votos := repositories.GetAllVotos()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(votos); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetVoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	participanteID, err := strconv.ParseInt(vars["participanteID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid participanteID format", http.StatusBadRequest)
		return
	}

	votacaoID, err := strconv.ParseInt(vars["votacaoID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid votacaoID format", http.StatusBadRequest)
		return
	}

	voto, exists := repositories.GetVotoByIDs(participanteID, votacaoID)
	if !exists {
		http.Error(w, "Voto not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(voto); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func CreateVoto(w http.ResponseWriter, r *http.Request) {
	var votoRequest struct {
		ParticipanteID int64 `json:"participanteId"`
		VotacaoID      int64 `json:"votacaoId"`
	}

	err := json.NewDecoder(r.Body).Decode(&votoRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Valida campos obrigatórios
	if votoRequest.ParticipanteID == 0 || votoRequest.VotacaoID == 0 {
		http.Error(w, "ParticipanteID and VotacaoID are required", http.StatusBadRequest)
		return
	}

	// Verifica se o participante existe
	participante, exists := repositories.GetParticipanteByID(votoRequest.ParticipanteID)
	if !exists {
		http.Error(w, "Participante not found", http.StatusNotFound)
		return
	}

	// Verifica se a votação existe
	votacao, exists := repositories.GetVotacaoByID(votoRequest.VotacaoID)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Cria voto
	voto := &entities.Voto{
		Participante: participante,
		Votacao:      votacao,
	}

	// Salva voto
	savedVoto := repositories.SaveVoto(voto)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(savedVoto); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
