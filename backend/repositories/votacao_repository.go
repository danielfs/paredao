package repositories

import (
	"database/sql"
	"log"

	"github.com/danielfs/paredao/backend/entities"
)

// GetAllVotacoes retrieves all votacoes from the database
func GetAllVotacoes() []*entities.Votacao {
	rows, err := DB.Query("SELECT id, descricao FROM votacoes")
	if err != nil {
		log.Printf("Error querying votacoes: %v", err)
		return []*entities.Votacao{}
	}
	defer rows.Close()

	votacoes := []*entities.Votacao{}
	for rows.Next() {
		v := &entities.Votacao{}
		if err := rows.Scan(&v.Id, &v.Descricao); err != nil {
			log.Printf("Error scanning votacao row: %v", err)
			continue
		}
		votacoes = append(votacoes, v)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating votacao rows: %v", err)
	}

	return votacoes
}

// GetVotacaoByID retrieves a votacao by ID
func GetVotacaoByID(id int64) (*entities.Votacao, bool) {
	v := &entities.Votacao{}
	err := DB.QueryRow("SELECT id, descricao FROM votacoes WHERE id = ?", id).
		Scan(&v.Id, &v.Descricao)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Error querying votacao by ID: %v", err)
		return nil, false
	}

	return v, true
}

// SaveVotacao saves a votacao to the database
func SaveVotacao(v *entities.Votacao) *entities.Votacao {
	if v.Id == 0 {
		// Insert new votacao
		result, err := DB.Exec(
			"INSERT INTO votacoes (descricao) VALUES (?)",
			v.Descricao,
		)
		if err != nil {
			log.Printf("Error inserting votacao: %v", err)
			return nil
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last insert ID: %v", err)
			return nil
		}

		v.Id = id
	} else {
		// Update existing votacao
		_, err := DB.Exec(
			"UPDATE votacoes SET descricao = ? WHERE id = ?",
			v.Descricao, v.Id,
		)
		if err != nil {
			log.Printf("Error updating votacao: %v", err)
			return nil
		}
	}

	return v
}

// DeleteVotacaoByID deletes a votacao by ID
func DeleteVotacaoByID(id int64) bool {
	result, err := DB.Exec("DELETE FROM votacoes WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting votacao: %v", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return false
	}

	return rowsAffected > 0
}

// AddParticipanteToVotacaoInDB adds a participante to a votacao
func AddParticipanteToVotacaoInDB(participanteID, votacaoID int64) bool {
	// Check if participante and votacao exist
	_, participanteExists := GetParticipanteByID(participanteID)
	_, votacaoExists := GetVotacaoByID(votacaoID)

	if !participanteExists || !votacaoExists {
		log.Printf("Cannot add participante to votacao: participante or votacao does not exist")
		return false
	}

	// Check if the relationship already exists
	var exists bool
	err := DB.QueryRow(
		"SELECT 1 FROM votacao_participante WHERE participante_id = ? AND votacao_id = ?",
		participanteID, votacaoID,
	).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking if participante is already in votacao: %v", err)
		return false
	}

	if err == nil {
		// Relationship already exists
		return true
	}

	// Insert new relationship
	_, err = DB.Exec(
		"INSERT INTO votacao_participante (participante_id, votacao_id) VALUES (?, ?)",
		participanteID, votacaoID,
	)

	if err != nil {
		log.Printf("Error adding participante to votacao: %v", err)
		return false
	}

	return true
}
