package token

type Token struct {
	name    string `db:"name" json:"name"`
	address string `db:"address" json:"address"`
}
