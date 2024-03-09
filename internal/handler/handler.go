package handler

import (
	"Kafka/configs"
	"Kafka/internal/models"
	"Kafka/internal/producer"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	configs *configs.Config
}

func NewCommentHandler(configs *configs.Config) *CommentHandler {
	return &CommentHandler{
		configs: configs,
	}
}

func (u *CommentHandler) Main() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("Welcome Kafka Application")
	}
}

func (c *CommentHandler) CreateComment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		cmt := models.Comment{}

		if err := ctx.BodyParser(&cmt); err != nil {
			log.Println(err)

			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		cmtInBytes, err := json.Marshal(cmt)
		if err != nil {
			log.Println(err)

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		err = producer.PushCommentToQueue(cmtInBytes, c.configs)
		if err != nil {
			log.Printf("Failed to send message %v", err)

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Comment pushed successfully",
			"comment": cmt,
		})
	}
}
