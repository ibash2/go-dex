package repository

import "go-dex/internal/pkg/token"

type Repository interface {
	GetTokens(tokens *[]token.Token) error
	AddUser(address string, inviterId int) error
}
