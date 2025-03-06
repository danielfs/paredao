package repositories

import (
	"database/sql"
	"log"

	"github.com/danielfs/paredao/backend/entities"
)

func GetAllParticipantes() []*entities.Participante {
	rows, err := DB.Query("SELECT id, nome, url_foto FROM participantes")
	if err != nil {
		log.Printf("Error querying participantes: %v", err)
		return []*entities.Participante{}
	}
	defer rows.Close()

	participantes := []*entities.Participante{}
	for rows.Next() {
		p := &entities.Participante{}
		if err := rows.Scan(&p.ID, &p.Nome, &p.URLFoto); err != nil {
			log.Printf("Error scanning participante row: %v", err)
			continue
		}
		participantes = append(participantes, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating participante rows: %v", err)
	}

	return participantes
}

func GetParticipanteByID(id int64) (*entities.Participante, bool) {
	p := &entities.Participante{}
	err := DB.QueryRow("SELECT id, nome, url_foto FROM participantes WHERE id = ?", id).
		Scan(&p.ID, &p.Nome, &p.URLFoto)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Error querying participante by ID: %v", err)
		return nil, false
	}

	return p, true
}

func SaveParticipante(p *entities.Participante) *entities.Participante {
	if p.ID == 0 {
		// Insere novo participante
		result, err := DB.Exec(
			"INSERT INTO participantes (nome, url_foto) VALUES (?, ?)",
			p.Nome, p.URLFoto,
		)
		if err != nil {
			log.Printf("Error inserting participante: %v", err)
			return nil
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last insert ID: %v", err)
			return nil
		}

		p.ID = id
	} else {
		// Atualiza participante existente
		_, err := DB.Exec(
			"UPDATE participantes SET nome = ?, url_foto = ? WHERE id = ?",
			p.Nome, p.URLFoto, p.ID,
		)
		if err != nil {
			log.Printf("Error updating participante: %v", err)
			return nil
		}
	}

	return p
}

func DeleteParticipanteByID(id int64) bool {
	result, err := DB.Exec("DELETE FROM participantes WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting participante: %v", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return false
	}

	return rowsAffected > 0
}

func GetParticipantesByVotacaoID(votacaoID int64) []*entities.Participante {
	query := `
		SELECT p.id, p.nome, p.url_foto
		FROM participantes p
		JOIN votacao_participante vp ON p.id = vp.participante_id
		WHERE vp.votacao_id = ?
	`

	rows, err := DB.Query(query, votacaoID)
	if err != nil {
		log.Printf("Error querying participantes by votacao ID: %v", err)
		return []*entities.Participante{}
	}
	defer rows.Close()

	participantes := []*entities.Participante{}
	for rows.Next() {
		p := &entities.Participante{}
		if err := rows.Scan(&p.ID, &p.Nome, &p.URLFoto); err != nil {
			log.Printf("Error scanning participante row: %v", err)
			continue
		}
		participantes = append(participantes, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating participante rows: %v", err)
	}

	return participantes
}
