package controllers

import (
	"demo-todo/models"
	"demo-todo/services"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// func CreateTodo(c *fiber.Ctx) error {
// 	var todo models.Todo
// 	if err := c.BodyParser(&todo); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request",
// 		})
// 	}

// 	conn := services.ConnectRabbitMQ()
// 	defer conn.Close()

// 	err := services.PublishMessage(conn, "todo_queue", todo.Title)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to publish message to RabbitMQ",
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "Todo created successfully",
// 		"todo":    todo,
// 	})
// }

func CreateTodos(c *fiber.Ctx) error {
	var todos []models.Todo
	if err := c.BodyParser(&todos); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	go func(todos []models.Todo) {
		conn := services.ConnectRabbitMQ()
		defer conn.Close()
		for _, todo := range todos {
			err := services.PublishMessage(conn, "todo_queue", todo.Title)
			if err != nil {
				log.Printf("Failed to publish todo: %v", err)
			}
			time.Sleep(1 * time.Second)
		}

	}(todos)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Todos received and will be processed",
	})
}
