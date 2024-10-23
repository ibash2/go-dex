package event

type BaseMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
