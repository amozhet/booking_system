package messaging

import (
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	if err := r.Channel.Close(); err != nil {
		log.Println("Failed to close RabbitMQ channel:", err)
	}
	if err := r.Conn.Close(); err != nil {
		log.Println("Failed to close RabbitMQ connection:", err)
	}
}
