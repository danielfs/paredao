package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/danielfs/paredao/backend/entities"
	"github.com/danielfs/paredao/backend/repositories"
)

func GetParticipantes(w http.ResponseWriter, r *http.Request) {
	participantes := repositories.GetAllParticipantes()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(participantes); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

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
	if err := json.NewEncoder(w).Encode(participante); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func CreateParticipante(w http.ResponseWriter, r *http.Request) {
	var participante entities.Participante
	err := json.NewDecoder(r.Body).Decode(&participante)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Valida campos obrigatórios
	if participante.Nome == "" {
		http.Error(w, "Nome is required", http.StatusBadRequest)
		return
	}

	// Salva participante
	savedParticipante := repositories.SaveParticipante(&participante)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(savedParticipante); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func UpdateParticipante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Verifica se o participante existe
	_, exists := repositories.GetParticipanteByID(id)
	if !exists {
		http.Error(w, "Participante not found", http.StatusNotFound)
		return
	}

	// Decodifica o corpo da requisição
	var participante entities.Participante
	err = json.NewDecoder(r.Body).Decode(&participante)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Garante que o ID corresponde ao parâmetro do caminho
	participante.ID = id

	// Valida campos obrigatórios
	if participante.Nome == "" {
		http.Error(w, "Nome is required", http.StatusBadRequest)
		return
	}

	// Salva o participante atualizado
	updatedParticipante := repositories.SaveParticipante(&participante)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedParticipante); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

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
