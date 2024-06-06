package messaging

import (
	"clientManage/internal/data"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type UserMessaging interface {
	PublishUserCreated(user *data.UserModel) error
	Close() error
}

type UserMessagingImpl struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewUserMessaging(rabbitMQUrl string) (*UserMessagingImpl, error) {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &UserMessagingImpl{connection: conn, channel: ch}, nil
}

func (m *UserMessagingImpl) PublishUserCreated(user *data.UserModel) error {
	err := m.channel.ExchangeDeclare(
		"user_exchange",
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

	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = m.channel.Publish(
		"user_exchange",
		"user.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish user created message: %v", err)
		return err
	}

	log.Printf("User created message published: %v", user)
	return nil
}

func (m *UserMessagingImpl) Close() error {
	if err := m.channel.Close(); err != nil {
		return err
	}
	return m.connection.Close()
}
