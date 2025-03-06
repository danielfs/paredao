package entities

type Participante struct {
	ID      int64  `json:"id"`
	Nome    string `json:"nome"`
	URLFoto string `json:"urlFoto"`
}
