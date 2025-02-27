package entities

type ParticipanteTotalResponse struct {
	ParticipanteId int64  `json:"participanteId"`
	Nome           string `json:"nome"`
	Total          int    `json:"total"`
}
