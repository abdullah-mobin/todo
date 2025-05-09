package routes

import (
	"demo-todo/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	todo := app.Group("/api/v1")

	todo.Post("/create", controllers.CreateTodos)

}
