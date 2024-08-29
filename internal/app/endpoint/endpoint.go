package endpoint

import (
	"encoding/json"
	"go-dex/internal/pkg/token"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service interface {
	GetTokens() []token.Token
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
	tokens, err := json.Marshal(e.service.GetTokens())

	if err != nil {
		return nil
	}

	err = ctx.JSONBlob(http.StatusOK, tokens)

	if err != nil {
		return err
	}
	return nil
}
