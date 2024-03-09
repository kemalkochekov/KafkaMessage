package server

import (
	"Kafka/configs"
	"Kafka/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(incomingRoutes *fiber.App, configs *configs.Config) {
	server := handler.NewCommentHandler(configs)
	api := incomingRoutes.Group("/api/v1")
	api.Post("/comments", server.CreateComment())
	api.Get("/", server.Main())
}
