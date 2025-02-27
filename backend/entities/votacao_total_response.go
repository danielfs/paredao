package entities

type VotacaoTotalResponse struct {
	VotacaoId int64 `json:"votacaoId"`
	Total     int   `json:"total"`
}
