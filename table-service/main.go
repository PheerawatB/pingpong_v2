package main

import (
	"fmt"

	"table-service/server"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Define routes for the Table service
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Service Table!")
	})

	app.Get("/ping-power", func(c *fiber.Ctx) error {
		fmt.Println("Received request at /ping-power")
		power := c.QueryInt("power")
		name := c.Query("name")
		if power <= 0 {
			return c.Status(fiber.StatusBadRequest).SendString("Power must be a positive integer")
		}
		if name == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Name is required")
		}
		lowPower := server.BallPowerTo(uint(power), name)
		return c.SendString(fmt.Sprintf("%d", lowPower))
	})

	if err := app.Listen(":8889"); err != nil {
		fmt.Println("Failed to start Table Service:", err)
	}

}
