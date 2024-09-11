package app

import (
	"fmt"
	"go-dex/internal/app/endpoint"
	"go-dex/internal/app/repository/mock"
	"go-dex/internal/app/service"

	"github.com/labstack/echo/v4"
)

type App struct {
	endpoint *endpoint.Endpoint
	service  *service.Service
	echo     *echo.Echo
}

func New() (*App, error) {
	app := &App{}

	db, err := mock.New()
	// db, err := sqlxx.New(sqlxx.LoadConfig())
	if err != nil {
		return nil, err
	}

	app.service = service.New(db)
	app.endpoint = endpoint.New(app.service)
	app.echo = echo.New()

	app.echo.GET("/tokens", app.endpoint.GetTokens)
	app.echo.POST("/user", app.endpoint.CreateUser)

	return app, nil
}

func (a *App) Run() error {
	fmt.Println("server running")

	err := a.echo.Start(":8080")
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
