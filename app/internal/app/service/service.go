package service

import (
	"fmt"

	"go-dex/internal/app/repository"
	"go-dex/internal/pkg/token"
)

type Service struct {
	db repository.Repository
}

func New(db repository.Repository) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetTokens() ([]token.Token, error) {
	var tokens []token.Token
	err := s.db.GetTokens(&tokens)

	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}

	return tokens, nil
}

func (s *Service) CreateUser(address string, inviterId int) error {
	err := s.db.AddUser(address, inviterId)

	if err != nil {
		fmt.Printf("failed to create user: %v", err)
	}
	return nil
}
