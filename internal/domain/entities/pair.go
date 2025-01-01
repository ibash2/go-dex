package entities

import "go-dex/internal/domain/events"

type Pair struct {
	BaseEntity
	ID     string `json:id`
	Token0 string `json:token0`
	Token1 string `json:token1`
}

func CreatePair(id string, token0 string, token1 string) *Pair {
	exmp := Pair{ID: id, Token0: token0, Token1: token1}
	exmp.AddEvent(events.NewPairEvent{PairID: id})
	return &exmp
}
