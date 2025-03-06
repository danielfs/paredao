package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/danielfs/paredao/backend/entities"
	"github.com/danielfs/paredao/backend/repositories"
)

func GetVotacoes(w http.ResponseWriter, r *http.Request) {
	votacoes := repositories.GetAllVotacoes()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(votacoes); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

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
	if err := json.NewEncoder(w).Encode(votacao); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func CreateVotacao(w http.ResponseWriter, r *http.Request) {
	var votacao entities.Votacao
	err := json.NewDecoder(r.Body).Decode(&votacao)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Valida campos obrigatórios
	if votacao.Descricao == "" {
		http.Error(w, "Descricao is required", http.StatusBadRequest)
		return
	}

	// Salva votação
	savedVotacao := repositories.SaveVotacao(&votacao)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(savedVotacao); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func UpdateVotacao(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Verifica se a votação existe
	_, exists := repositories.GetVotacaoByID(id)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Decodifica o corpo da requisição
	var votacao entities.Votacao
	err = json.NewDecoder(r.Body).Decode(&votacao)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Garante que o ID corresponde ao parâmetro do caminho
	votacao.ID = id

	// Valida campos obrigatórios
	if votacao.Descricao == "" {
		http.Error(w, "Descricao is required", http.StatusBadRequest)
		return
	}

	// Salva a votação atualizada
	updatedVotacao := repositories.SaveVotacao(&votacao)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedVotacao); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

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

func GetVotacaoParticipantes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Verifica se a votação existe
	_, exists := repositories.GetVotacaoByID(id)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	participantes := repositories.GetParticipantesByVotacaoID(id)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(participantes); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func AddParticipanteToVotacao(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	votacaoID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid votacao ID format", http.StatusBadRequest)
		return
	}

	// Verifica se a votação existe
	_, exists := repositories.GetVotacaoByID(votacaoID)
	if !exists {
		http.Error(w, "Votacao not found", http.StatusNotFound)
		return
	}

	// Analisa o corpo da requisição
	var request struct {
		ParticipanteID int64 `json:"participanteId"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Valida campos obrigatórios
	if request.ParticipanteID == 0 {
		http.Error(w, "ParticipanteId is required", http.StatusBadRequest)
		return
	}

	// Verifica se o participante existe
	participante, exists := repositories.GetParticipanteByID(request.ParticipanteID)
	if !exists {
		http.Error(w, "Participante not found", http.StatusNotFound)
		return
	}

	// Adiciona participante à votação
	success := repositories.AddParticipanteToVotacaoInDB(request.ParticipanteID, votacaoID)
	if !success {
		http.Error(w, "Failed to add participante to votacao", http.StatusInternalServerError)
		return
	}

	// Retorna o participante que foi adicionado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(participante); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
