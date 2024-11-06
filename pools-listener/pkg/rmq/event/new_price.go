package event

type NewPriceEvent struct {
	Address string  `json:"address"`
	Price   float64 `json:"price"`
}
