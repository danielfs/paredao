package repositories

import (
	"github.com/danielfs/paredao/backend/entities"
)

func GetTotalVotesForVotacao(votacaoID int64) (int, error) {
	var total int
	query := "SELECT COUNT(*) FROM votos WHERE votacao_id = ?"

	err := DB.QueryRow(query, votacaoID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func GetTotalVotesByParticipante(votacaoID int64) ([]entities.ParticipanteTotalResponse, error) {
	// Primeiro, obtém todos os participantes para esta votação
	participants := GetParticipantesByVotacaoID(votacaoID)

	// Cria um mapa para armazenar os totais de votos para cada participante
	participantTotals := make(map[int64]int)

	// Consulta para obter contagens de votos para participantes que têm votos
	query := `
		SELECT p.id, COUNT(*) as total
		FROM votos v
		JOIN participantes p ON v.participante_id = p.id
		WHERE v.votacao_id = ?
		GROUP BY p.id
	`

	rows, err := DB.Query(query, votacaoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Preenche o mapa com contagens de votos
	for rows.Next() {
		var participanteID int64
		var total int
		if err := rows.Scan(&participanteID, &total); err != nil {
			return nil, err
		}
		participantTotals[participanteID] = total
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cria resposta com todos os participantes, incluindo aqueles com zero votos
	totals := make([]entities.ParticipanteTotalResponse, 0, len(participants))
	for _, p := range participants {
		total := participantTotals[p.ID] // Will be 0 if no votes
		totals = append(totals, entities.ParticipanteTotalResponse{
			ParticipanteID: p.ID,
			Nome:           p.Nome,
			Total:          total,
		})
	}

	return totals, nil
}

func GetTotalVotesByHour(votacaoID int64) ([]entities.HourlyTotalResponse, error) {
	query := `
		SELECT HOUR(data_hora) as hour, COUNT(*) as total
		FROM votos
		WHERE votacao_id = ?
		GROUP BY HOUR(data_hora)
		ORDER BY hour
	`

	rows, err := DB.Query(query, votacaoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totals []entities.HourlyTotalResponse
	for rows.Next() {
		var total entities.HourlyTotalResponse
		if err := rows.Scan(&total.Hour, &total.Total); err != nil {
			return nil, err
		}
		totals = append(totals, total)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Inicializa todas as 24 horas com contagens zero
	hourlyTotals := make([]entities.HourlyTotalResponse, 24)
	for i := 0; i < 24; i++ {
		hourlyTotals[i] = entities.HourlyTotalResponse{Hour: i, Total: 0}
	}

	// Atualiza com contagens reais
	for _, total := range totals {
		hourlyTotals[total.Hour] = total
	}

	return hourlyTotals, nil
}
