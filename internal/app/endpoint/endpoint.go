package endpoint

import (
	"encoding/json"
	"go-dex/internal/pkg/token"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUser struct {
	Address string `json:"address"`
}

type Service interface {
	GetTokens() []token.Token
	CreateUser(address string) error
}

type Endpoint struct {
	service Service
}

func New(service Service) *Endpoint {
	return &Endpoint{
		service: service,
	}
}

func (e *Endpoint) GetTokens(ctx echo.Context) error {
	tokens, err := json.Marshal(token.Token{Name: "test", Address: "niger"})

	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err = ctx.JSONBlob(http.StatusOK, tokens)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return nil
}

func (e *Endpoint) CreateUser(ctx echo.Context) error {
	user := new(CreateUser)
	if err := ctx.Bind(user); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err := e.service.CreateUser(user.Address)

	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, user)
}
