package entities

type ParticipanteTotalResponse struct {
	ParticipanteID int64  `json:"participanteId"`
	Nome           string `json:"nome"`
	Total          int    `json:"total"`
}
