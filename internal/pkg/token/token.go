package token

type Token struct {
	Symbol  string `db:"symbol" json:"symbol"`
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
}
