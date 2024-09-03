package service

import (
	"fmt"
	"log"

	"go-dex/internal/app/sqlxx"
	"go-dex/internal/pkg/token"
)

type Service struct {
	db *sqlxx.DB
}

func New(db *sqlxx.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetTokens() []token.Token {
	var token []token.Token

	err := s.db.Select(&token, "SELECT symbol, name, address FROM token")

	if err != nil {
		fmt.Printf("failed to get tokens: %v", err)
	}

	return token
}

func (s *Service) CreateUser(address string) error {
	result, err := s.db.Exec(
		`INSERT INTO "user" (address, points) VALUES ($1, $2)`,
		address, 100)

	if err != nil {
		log.Fatalf("Error inserting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error getting affected rows: %v", err)
	}

	fmt.Printf("Number of rows affected: %d\n", rowsAffected)

	return nil
}
