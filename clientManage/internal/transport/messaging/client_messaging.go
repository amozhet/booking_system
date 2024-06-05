package messaging

import (
	"encoding/json"
	"log"

	"clientManage/internal/domain/model"
	"github.com/streadway/amqp"
)

type ClientMessaging interface {
	PublishClientCreated(client *model.Client) error
	Close() error
}

type ClientMessagingImpl struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewClientMessaging(rabbitMQUrl string) (*ClientMessagingImpl, error) {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &ClientMessagingImpl{connection: conn, channel: ch}, nil
}

func (m *ClientMessagingImpl) PublishClientCreated(client *model.Client) error {
	err := m.channel.ExchangeDeclare(
		"client_exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(client)
	if err != nil {
		return err
	}

	err = m.channel.Publish(
		"client_exchange",
		"client.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish client created message: %v", err)
		return err
	}

	log.Printf("Client created message published: %v", client)
	return nil
}

func (m *ClientMessagingImpl) Close() error {
	if err := m.channel.Close(); err != nil {
		return err
	}
	return m.connection.Close()
}
