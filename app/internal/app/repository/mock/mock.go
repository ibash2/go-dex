package mock

import (
	"go-dex/internal/pkg/token"
	"go-dex/internal/pkg/user"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
)

// DBrepo is a wrapper around sqlx.DBrepo
type mockRepo struct {
	tokens []token.Token
	users  []user.User
}

// GetUsers implements repository.Repository.
func (db *mockRepo) GetUsers() ([]string, error) {
	panic("unimplemented")
}

// New - create new connect to DB
func New() (db *mockRepo, myerr error) {
	// Создаем новый сервис
	db = &mockRepo{
		tokens: []token.Token{
			{Symbol: "test", Name: "test", Address: "test"},
		},
		users: []user.User{
			{Id: 1, Address: "test", InviterId: 1923, Points: 1}},
	}

	return db, nil
}

func (db *mockRepo) GetTokens(tokens *[]token.Token) error {
	*tokens = db.tokens
	return nil
}

func (db *mockRepo) AddUser(address string, inviterId int) error {
	for _, users := range db.users {
		if users.Address == address {
			return nil
		} else if inviterId != 0 && inviterId == users.InviterId {

			users.Points += 100
			db.users = append(db.users, user.User{
				Id:        len(db.users) + 1,
				Address:   address,
				InviterId: inviterId,
				Points:    200,
			})

			return nil
		}
	}
	db.users = append(db.users, user.User{
		Id:        len(db.users) + 1,
		Address:   address,
		InviterId: 192,
	})
	return nil
}

func (db *mockRepo) Get() error {
	return nil
}
func (db *mockRepo) Post() error {
	return nil
}
