package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

type RoomMessagingClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRoomMessagingClient(rabbitMQUrl string) (*RoomMessagingClient, error) {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"room_updates",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RoomMessagingClient{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (c *RoomMessagingClient) PublishUpdate(message string) error {
	return c.channel.Publish(
		"",
		c.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (c *RoomMessagingClient) Close() {
	if err := c.channel.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}
