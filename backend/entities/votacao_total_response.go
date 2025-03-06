package entities

type VotacaoTotalResponse struct {
	VotacaoID int64 `json:"votacaoId"`
	Total     int   `json:"total"`
}
