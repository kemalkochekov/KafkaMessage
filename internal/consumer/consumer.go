package consumer

import (
	"Kafka/configs"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
}

func NewConsumer(configs *configs.Config) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(configs.Kafka.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to open kafka consumer")
	}

	return &Consumer{
		consumer: consumer,
	}, nil
}
func (c *Consumer) ReadMessages(topic string) error {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Failed to open partitions %w", err)
	}

	msgCount := 0
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	doneCh := make(chan struct{})

	for _, partition := range partitions {
		go func(partition int32) {
			partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			if err != nil {
				return
			}

			defer partitionConsumer.Close()

			for {
				select {
				case err = <-partitionConsumer.Errors():
					fmt.Println(err)
				case msg := <-partitionConsumer.Messages():
					msgCount++
					fmt.Printf("Recieved message Count: %d, Topic: %s, Message: %s\n", msgCount, msg.Topic, string(msg.Value))

				case <-sigChan:
					log.Println("Shutting down consumer for partition", partition)
					doneCh <- struct{}{}
				}
			}
		}(partition)
	}

	<-doneCh
	fmt.Println("Processed ", msgCount, " messages")

	defer func() {
		err = c.consumer.Close()
		if err != nil {
			log.Printf("Failed to close consumer %v", err)
		}
	}()

	return fmt.Errorf("Consumer canceled %w", err)
}
