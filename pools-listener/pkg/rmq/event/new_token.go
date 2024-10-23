package event

type NewTokenEvent struct {
	Address  string `json:"address"`
	Address2 string `json:"address2"`
	Chain    string `json:"chain"`
}
