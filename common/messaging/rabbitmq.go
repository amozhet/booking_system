package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

func ConnectRabbitMQ(rabbitMQUrl string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	return conn, nil
}
