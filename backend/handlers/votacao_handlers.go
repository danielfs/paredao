package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danielfs/paredao/backend/entities"
	"github.com/danielfs/paredao/backend/repositories"
	"github.com/gorilla/mux"
)

// GetVotacoes handles GET /votacoes
func GetVotacoes(w http.ResponseWriter, r *http.Request) {
	votacoes := repositories.GetAllVotacoes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votacoes)
}

// GetVotacao handles GET /votacoes/{id}
func GetVotacao(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	votacao, exists := repositories.GetVotacaoByID(id)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votacao)
}

// CreateVotacao handles POST /votacoes
func CreateVotacao(w http.ResponseWriter, r *http.Request) {
	var votacao entities.Votacao
	err := json.NewDecoder(r.Body).Decode(&votacao)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if votacao.Descricao == "" {
		http.Error(w, "Descricao is required", http.StatusBadRequest)
		return
	}

	// Save votacao
	savedVotacao := repositories.SaveVotacao(&votacao)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedVotacao)
}

// UpdateVotacao handles PUT /votacoes/{id}
func UpdateVotacao(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Check if votacao exists
	_, exists := repositories.GetVotacaoByID(id)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Decode request body
	var votacao entities.Votacao
	err = json.NewDecoder(r.Body).Decode(&votacao)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches path parameter
	votacao.Id = id

	// Validate required fields
	if votacao.Descricao == "" {
		http.Error(w, "Descricao is required", http.StatusBadRequest)
		return
	}

	// Save updated votacao
	updatedVotacao := repositories.SaveVotacao(&votacao)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedVotacao)
}

// DeleteVotacao handles DELETE /votacoes/{id}
func DeleteVotacao(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	success := repositories.DeleteVotacaoByID(id)
	if !success {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetVotacaoParticipantes handles GET /votacoes/{id}/participantes
func GetVotacaoParticipantes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Check if votacao exists
	_, exists := repositories.GetVotacaoByID(id)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	participantes := repositories.GetParticipantesByVotacaoID(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participantes)
}

// AddParticipanteToVotacao handles POST /votacoes/{id}/participantes
func AddParticipanteToVotacao(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	votacaoID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid votacao ID format", http.StatusBadRequest)
		return
	}

	// Check if votacao exists
	_, exists := repositories.GetVotacaoByID(votacaoID)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Parse request body
	var request struct {
		ParticipanteID int64 `json:"participanteId"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.ParticipanteID == 0 {
		http.Error(w, "ParticipanteId is required", http.StatusBadRequest)
		return
	}

	// Check if participante exists
	participante, exists := repositories.GetParticipanteByID(request.ParticipanteID)
	if !exists {
		http.Error(w, "Participante not found", http.StatusNotFound)
		return
	}

	// Add participante to votacao
	success := repositories.AddParticipanteToVotacaoInDB(request.ParticipanteID, votacaoID)
	if !success {
		http.Error(w, "Failed to add participante to votacao", http.StatusInternalServerError)
		return
	}

	// Return the participante that was added
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(participante)
}
