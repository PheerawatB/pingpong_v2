package main

import (
	"context"
	"fmt"
	"strconv"

	"player-service/server"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// initialize variables
var mongoClient *mongo.Client

func main() {
	var err error
	mongoURI := "mongodb://127.0.0.1:27017/?serverSelectionTimeoutMS=5000"
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	defer mongoClient.Disconnect(context.TODO())
	server.CountMatch, _ = server.GetLastMatchID(mongoClient)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Player Service")
	})

	app.Post("/new-match", func(c *fiber.Ctx) error {
		_ = server.NewMatch()
		server.LogMatchResultToMongoDB(server.CountMatch, server.LogMatch, mongoClient)
		return c.SendString(server.LogMatch)
	})

	app.Get("/match", func(c *fiber.Ctx) error {
		listMatch, _ := server.GetAllMatches(mongoClient)
		return c.JSON(listMatch)
	})

	app.Get("/match/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid match ID")
		}

		res, err := server.GetMatchID(mongoClient, uint(id))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to get match data")
		}

		return c.JSON(res)
	})

	if err := app.Listen(":8888"); err != nil {
		fmt.Println("Failed to start Player Service:", err)
	}
}
