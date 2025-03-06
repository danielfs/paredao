package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/danielfs/paredao/backend/entities"
)

func GetAllVotos() []*entities.Voto {
	query := `
		SELECT v.participante_id, v.votacao_id, v.data_hora,
			   p.id, p.nome, p.url_foto,
			   vt.id, vt.descricao
		FROM votos v
		JOIN participantes p ON v.participante_id = p.id
		JOIN votacoes vt ON v.votacao_id = vt.id
	`

	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error querying votos: %v", err)
		return []*entities.Voto{}
	}
	defer rows.Close()

	votos := []*entities.Voto{}
	for rows.Next() {
		v := &entities.Voto{
			Participante: &entities.Participante{},
			Votacao:      &entities.Votacao{},
		}

		if err := rows.Scan(
			&v.Participante.ID, &v.Votacao.ID, &v.DataHora,
			&v.Participante.ID, &v.Participante.Nome, &v.Participante.URLFoto,
			&v.Votacao.ID, &v.Votacao.Descricao,
		); err != nil {
			log.Printf("Error scanning voto row: %v", err)
			continue
		}

		votos = append(votos, v)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating voto rows: %v", err)
	}

	return votos
}

func GetVotoByIDs(participanteID, votacaoID int64) (*entities.Voto, bool) {
	query := `
		SELECT v.participante_id, v.votacao_id, v.data_hora,
			   p.id, p.nome, p.url_foto,
			   vt.id, vt.descricao
		FROM votos v
		JOIN participantes p ON v.participante_id = p.id
		JOIN votacoes vt ON v.votacao_id = vt.id
		WHERE v.participante_id = ? AND v.votacao_id = ?
	`

	v := &entities.Voto{
		Participante: &entities.Participante{},
		Votacao:      &entities.Votacao{},
	}

	err := DB.QueryRow(query, participanteID, votacaoID).Scan(
		&v.Participante.ID, &v.Votacao.ID, &v.DataHora,
		&v.Participante.ID, &v.Participante.Nome, &v.Participante.URLFoto,
		&v.Votacao.ID, &v.Votacao.Descricao,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Error querying voto by IDs: %v", err)
		return nil, false
	}

	return v, true
}

func SaveVoto(v *entities.Voto) *entities.Voto {
	// Verifica se o participante e a votação existem
	_, participanteExists := GetParticipanteByID(v.Participante.ID)
	_, votacaoExists := GetVotacaoByID(v.Votacao.ID)

	if !participanteExists || !votacaoExists {
		log.Printf("Cannot save voto: participante or votacao does not exist")
		return nil
	}

	// Define o timestamp se não fornecido
	if v.DataHora.IsZero() {
		v.DataHora = time.Now()
	}

	// Insere novo voto
	_, err := DB.Exec(
		"INSERT INTO votos (participante_id, votacao_id, data_hora) VALUES (?, ?, ?)",
		v.Participante.ID, v.Votacao.ID, v.DataHora,
	)
	if err != nil {
		log.Printf("Error inserting voto: %v", err)
		return nil
	}

	return v
}
