package entities

import "time"

type Voto struct {
	Participante *Participante `json:"participante"`
	Votacao      *Votacao      `json:"votacao"`
	DataHora     time.Time     `json:"dataHora"`
}
