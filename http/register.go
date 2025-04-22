package http

import (
	"medods/http/handlers/tasks"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func Register(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New())

	tasks.RegisterTasks(app)
}
