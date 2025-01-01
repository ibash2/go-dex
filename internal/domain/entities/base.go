package entities

type BaseEntity struct {
	Events []interface{} `json:"events"`
}

func (b *BaseEntity) AddEvent(event interface{}) {
	b.Events = append(b.Events, event)
}
