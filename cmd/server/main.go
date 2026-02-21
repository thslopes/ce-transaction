package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/thslopes/ce-transactions/internal/app"
)

func main() {
	// Initialize a new Fiber app
	app := app.App{
		FiberApp: fiber.New(),
	}

	// Start the server on port 3000
	log.Fatal(app.FiberApp.Listen(":3000"))
}
