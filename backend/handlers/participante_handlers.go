package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danielfs/paredao/backend/entities"
	"github.com/danielfs/paredao/backend/repositories"
	"github.com/gorilla/mux"
)

// GetParticipantes handles GET /participantes
func GetParticipantes(w http.ResponseWriter, r *http.Request) {
	participantes := repositories.GetAllParticipantes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participantes)
}

// GetParticipante handles GET /participantes/{id}
func GetParticipante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	participante, exists := repositories.GetParticipanteByID(id)
	if !exists {
		http.Error(w, "Participante not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participante)
}

// CreateParticipante handles POST /participantes
func CreateParticipante(w http.ResponseWriter, r *http.Request) {
	var participante entities.Participante
	err := json.NewDecoder(r.Body).Decode(&participante)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if participante.Nome == "" {
		http.Error(w, "Nome is required", http.StatusBadRequest)
		return
	}

	// Save participante
	savedParticipante := repositories.SaveParticipante(&participante)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedParticipante)
}

// UpdateParticipante handles PUT /participantes/{id}
func UpdateParticipante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Check if participante exists
	_, exists := repositories.GetParticipanteByID(id)
	if !exists {
		http.Error(w, "Participante not found", http.StatusNotFound)
		return
	}

	// Decode request body
	var participante entities.Participante
	err = json.NewDecoder(r.Body).Decode(&participante)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure ID matches path parameter
	participante.Id = id

	// Validate required fields
	if participante.Nome == "" {
		http.Error(w, "Nome is required", http.StatusBadRequest)
		return
	}

	// Save updated participante
	updatedParticipante := repositories.SaveParticipante(&participante)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedParticipante)
}

// DeleteParticipante handles DELETE /participantes/{id}
func DeleteParticipante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	success := repositories.DeleteParticipanteByID(id)
	if !success {
		http.Error(w, "Participante not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
