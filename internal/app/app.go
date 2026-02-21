package app

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type App struct {
	FiberApp               *fiber.App
	GetTransactionsHandler func(c fiber.Ctx) error
}

func MockApp() *App {
	return &App{
		FiberApp: fiber.New(),
		GetTransactionsHandler: func(c fiber.Ctx) error {
			return c.SendString(fmt.Sprintf("Mock Get Transactions Handler - Method: %s", c.Method()))
		},
	}
}
