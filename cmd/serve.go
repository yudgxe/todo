/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"medods/database"
	"medods/database/dao"
	httpreg "medods/http"
	"medods/logger"
	"medods/utils"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start service",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger(logrus.Level(viper.GetInt("logger.level")))

		url := fmt.Sprintf("postgres://%v:%v@%v/%v",
			viper.GetString("db.user"),
			viper.GetString("db.password"),
			viper.GetString("db.address"),
			viper.GetString("db.database"),
		)

		database.MustInitDatabase(context.Background(), url)

		app := fiber.New(fiber.Config{
			ErrorHandler: func(ctx fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError

				var fiberErr *fiber.Error
				if errors.As(err, &fiberErr) {
					code = fiberErr.Code
				}

				var NotFoundErr *dao.ErrTaskNotFound
				if errors.As(err, &NotFoundErr) {
					code = http.StatusBadRequest
				}

				ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

				return ctx.Status(code).JSON(fiber.Map{"error": err.Error()})
			},
		})

		httpreg.Register(app)

		if err := app.Listen(viper.GetString("server.address")); err != nil {
			utils.Panicf("error on start server - %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
