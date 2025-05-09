package main

import (
	"context"
	"demo-todo/routes"
	"demo-todo/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	rabbitConn := services.ConnectRabbitMQ()
	defer rabbitConn.Close()

	_, err := services.CreateQueue(rabbitConn, "todo_queue")
	if err != nil {
		panic("Failed to create RabbitMQ queue")
	}

	mongoClient, err := services.ConnectMongoDB("mongodb://localhost:27017")
	if err != nil {
		panic("Failed to connect to MongoDB")
	}
	defer mongoClient.Disconnect(context.Background())

	go services.ConsumeMessages(rabbitConn, mongoClient, "todo_db", "todos", "todo_queue")

	app := fiber.New()
	routes.SetupRoutes(app)

	log.Println("Starting server on :3000")
	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
