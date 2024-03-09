package main

import (
	"Kafka/configs"
	"Kafka/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {
	router := fiber.New()

	router.Use(func(c *fiber.Ctx) error {
		// Log information about the incoming request
		println("Method:", c.Method(), "Path:", c.Path())
		return c.Next() // Move to the next middleware/handler
	})

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to get configuration data %v", err)
	}
	// Create Kafka Producer
	server.ApiRoutes(router, config)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(quit)

	go func() {
		<-quit
		// Perform graceful shut down
		os.Exit(0)
	}()

	err = router.Listen(config.Server.Host)
	if err != nil {
		return
	}
}
