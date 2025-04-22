package tasks

import (
	"medods/database/dao"
	"medods/database/model"
	"medods/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

func RegisterTasks(app *fiber.App) {
	app.Get("/tasks", get())
	app.Post("/tasks", create())
	app.Delete("/tasks/:id", delete())
	app.Put("/tasks/:id", update())
}

func get() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		tasks, err := dao.Tasks().List(ctx.Context())
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(tasks)
	}
}

func create() fiber.Handler {
	type body struct {
		Title       string           `json:"title"`
		Description *string          `json:"description"`
		Status      model.EnumStatus `json:"status"`
	}

	return func(ctx fiber.Ctx) error {
		var body body
		if err := ctx.Bind().Body(&body); err != nil {
			return err
		}

		now := time.Now()
		
		task := &model.Task{
			Title:       body.Title,
			Description: body.Description,
			Status:      utils.Ternary(body.Status == "", model.EnumStatusNew, body.Status),
			CreatedAt:   &now,
		}

		if err := dao.Tasks().Create(ctx.Context(), task); err != nil {
			return err
		}

		return ctx.Status(http.StatusCreated).JSON(task)
	}
}

func delete() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		id, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}

		return dao.Tasks().Delete(ctx.Context(), int32(id))
	}
}

func update() fiber.Handler {
	type body struct {
		Title       string           `json:"title"`
		Description *string          `json:"description"`
		Status      model.EnumStatus `json:"status"`
	}

	return func(ctx fiber.Ctx) error {
		id, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}

		var body body
		if err := ctx.Bind().Body(&body); err != nil {
			return err
		}

		if body.Status == "" {
			return fiber.NewError(http.StatusBadRequest, "status can't be empty")
		}

		now := time.Now()

		task := &model.Task{
			Id:          int32(id),
			Status:      body.Status,
			Title:       body.Title,
			Description: body.Description,
			UpdatedAt:   &now,
		}

		if err := dao.Tasks().Update(ctx.Context(), task); err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(task)
	}
}
