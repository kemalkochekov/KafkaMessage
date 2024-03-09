package main

import (
	"Kafka/configs"
	"Kafka/internal/consumer"
	"log"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to get configuration data %v", err)
	}

	worker, err := consumer.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer %v", err)
	}

	err = worker.ReadMessages(config.Kafka.Topic)
	if err != nil {
		log.Printf("Consumer closed %v", err)
	}
}
