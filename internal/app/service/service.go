package service

import (
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
	err := s.db.Select(&token, "SELECT * FROM token")
	if err != nil {
		return nil
	}

	return token
}

// var tokens []map[string]string
// 	for _, v := range token {
// 		tokens = append(tokens, map[string]string{
// 			"name":    v.Name,
// 			"address": v.Address,
// 		})
// 	}
// 	return tokens
