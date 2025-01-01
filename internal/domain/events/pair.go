package events

type NewPairEvent struct {
	BaseEvent
	PairID string `json:"pair_id"`
}
