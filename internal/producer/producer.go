package producer

import (
	"Kafka/configs"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func connectProducer(configs *configs.Config) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	// waits message delivered
	config.Producer.Return.Successes = true
	// retying send message max 5 times
	config.Producer.Retry.Max = 5
	// wants a message from broker leader replica acks = 1
	config.Producer.RequiredAcks = sarama.WaitForLocal

	syncProducer, err := sarama.NewSyncProducer(configs.Kafka.Brokers, config)
	if err != nil {
		return nil, err
	}

	return syncProducer, nil
}

func PushCommentToQueue(message []byte, configs *configs.Config) error {
	producer, err := connectProducer(configs)
	if err != nil {
		return err
	}

	defer func() {
		err = producer.Close()
		if err != nil {
			log.Printf("Failed to close producer %v", err)
		}
	}()

	msg := &sarama.ProducerMessage{
		Topic: configs.Kafka.Topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Failed to send message")
	}

	fmt.Printf("Message is stored in topic %s, partition %d, offset %d\n", configs.Kafka.Topic, partition, offset)

	return nil
}
