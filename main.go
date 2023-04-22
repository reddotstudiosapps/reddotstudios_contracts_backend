package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New()

	// Register the logger middleware
	app.Use(logger.New())

	// Define a basic GET route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	// Start the server on port 8080
	port := 8080
	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error starting server on port %d: %v", port, err)
	}
}
