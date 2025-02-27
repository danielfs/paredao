package repositories

import (
	"github.com/danielfs/paredao/backend/entities"
)

// GetTotalVotesForVotacao returns the total number of votes for a votacao
func GetTotalVotesForVotacao(votacaoId int64) (int, error) {
	var total int
	query := "SELECT COUNT(*) FROM votos WHERE votacao_id = ?"

	err := DB.QueryRow(query, votacaoId).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetTotalVotesByParticipante returns the total number of votes by participante for a votacao
func GetTotalVotesByParticipante(votacaoId int64) ([]entities.ParticipanteTotalResponse, error) {
	// First, get all participants for this votacao
	participants := GetParticipantesByVotacaoID(votacaoId)

	// Create a map to store vote totals for each participant
	participantTotals := make(map[int64]int)

	// Query to get vote counts for participants who have votes
	query := `
		SELECT p.id, COUNT(*) as total
		FROM votos v
		JOIN participantes p ON v.participante_id = p.id
		WHERE v.votacao_id = ?
		GROUP BY p.id
	`

	rows, err := DB.Query(query, votacaoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Fill the map with vote counts
	for rows.Next() {
		var participanteId int64
		var total int
		if err := rows.Scan(&participanteId, &total); err != nil {
			return nil, err
		}
		participantTotals[participanteId] = total
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Create response with all participants, including those with zero votes
	var totals []entities.ParticipanteTotalResponse
	for _, p := range participants {
		total := participantTotals[p.Id] // Will be 0 if no votes
		totals = append(totals, entities.ParticipanteTotalResponse{
			ParticipanteId: p.Id,
			Nome:           p.Nome,
			Total:          total,
		})
	}

	// Sort by total votes (descending)
	for i := 0; i < len(totals)-1; i++ {
		for j := i + 1; j < len(totals); j++ {
			if totals[i].Total < totals[j].Total {
				totals[i], totals[j] = totals[j], totals[i]
			}
		}
	}

	return totals, nil
}

// GetTotalVotesByHour returns the total number of votes per hour for a votacao
func GetTotalVotesByHour(votacaoId int64) ([]entities.HourlyTotalResponse, error) {
	query := `
		SELECT HOUR(data_hora) as hour, COUNT(*) as total
		FROM votos
		WHERE votacao_id = ?
		GROUP BY HOUR(data_hora)
		ORDER BY hour
	`

	rows, err := DB.Query(query, votacaoId)
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

	// Initialize all 24 hours with zero counts
	hourlyTotals := make([]entities.HourlyTotalResponse, 24)
	for i := 0; i < 24; i++ {
		hourlyTotals[i] = entities.HourlyTotalResponse{Hour: i, Total: 0}
	}

	// Update with actual counts
	for _, total := range totals {
		hourlyTotals[total.Hour] = total
	}

	return hourlyTotals, nil
}
