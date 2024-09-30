package endpoint

import (
	"net/http"
	"os/user"

	"github.com/labstack/echo/v4"
)

type CreateUser struct {
	Address string `json:"address"`
	Inviter int    `json:"inviter"`
}

type Service interface {
	GetTokens() ([]user.User, error)
	CreateUser(address string, inviterId int) error
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
	tokens, err := e.service.GetTokens()

	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	err = ctx.JSON(http.StatusOK, tokens)
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

	err := e.service.CreateUser(user.Address, user.Inviter)

	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, user)
}
