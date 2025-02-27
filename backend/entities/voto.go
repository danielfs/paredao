package entities

import "time"

type Voto struct {
	Participante *Participante
	Votacao      *Votacao
	DataHora     time.Time
}
