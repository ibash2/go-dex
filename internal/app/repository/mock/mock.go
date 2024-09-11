package mock

import (
	"go-dex/internal/pkg/token"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
)

// DBrepo is a wrapper around sqlx.DBrepo
type mockRepo struct {
	tokens []token.Token
	users  []string
}

// New - create new connect to DB
func New() (db *mockRepo, myerr error) {
	// Создаем новый сервис
	db = &mockRepo{
		tokens: []token.Token{
			{Symbol: "test", Name: "test", Address: "test"},
		},
		users: []string{},
	}

	return db, nil
}

func (db *mockRepo) GetTokens(tokens *[]token.Token) error {
	*tokens = db.tokens
	return nil
}

func (db *mockRepo) AddUser(address string, inviterId int) error {
	for _, user := range db.users {
		if user == address {
			return nil
		}
	}
	db.users = append(db.users, address)
	return nil
}
