package token

type Token struct {
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
}
