package repositories

import (
	"database/sql"
	"log"

	"github.com/danielfs/paredao/backend/entities"
)

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
		if err := rows.Scan(&v.ID, &v.Descricao); err != nil {
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

func GetVotacaoByID(id int64) (*entities.Votacao, bool) {
	v := &entities.Votacao{}
	err := DB.QueryRow("SELECT id, descricao FROM votacoes WHERE id = ?", id).
		Scan(&v.ID, &v.Descricao)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Error querying votacao by ID: %v", err)
		return nil, false
	}

	return v, true
}

func SaveVotacao(v *entities.Votacao) *entities.Votacao {
	if v.ID == 0 {
		// Insere nova votação
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

		v.ID = id
	} else {
		// Atualiza votação existente
		_, err := DB.Exec(
			"UPDATE votacoes SET descricao = ? WHERE id = ?",
			v.Descricao, v.ID,
		)
		if err != nil {
			log.Printf("Error updating votacao: %v", err)
			return nil
		}
	}

	return v
}

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

func AddParticipanteToVotacaoInDB(participanteID, votacaoID int64) bool {
	// Verifica se o participante e a votação existem
	_, participanteExists := GetParticipanteByID(participanteID)
	_, votacaoExists := GetVotacaoByID(votacaoID)

	if !participanteExists || !votacaoExists {
		log.Printf("Cannot add participante to votacao: participante or votacao does not exist")
		return false
	}

	// Verifica se o relacionamento já existe
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
		// Relacionamento já existe
		return true
	}

	// Insere novo relacionamento
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
