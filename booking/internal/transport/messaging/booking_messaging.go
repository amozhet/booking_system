package messaging

import (
	"encoding/json"
	"log"

	"booking/internal/domain/model"
	"github.com/streadway/amqp"
)

type BookingMessaging interface {
	PublishBookingCreated(booking *model.Booking) error
	Close() error
}

type BookingMessagingImpl struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewBookingMessaging(rabbitMQUrl string) (*BookingMessagingImpl, error) {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &BookingMessagingImpl{connection: conn, channel: ch}, nil
}

func (m *BookingMessagingImpl) PublishBookingCreated(booking *model.Booking) error {
	err := m.channel.ExchangeDeclare(
		"booking_exchange",
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

	body, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	err = m.channel.Publish(
		"booking_exchange",
		"booking.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish booking created message: %v", err)
		return err
	}

	log.Printf("Booking created message published: %v", booking)
	return nil
}

func (m *BookingMessagingImpl) Close() error {
	if err := m.channel.Close(); err != nil {
		return err
	}
	return m.connection.Close()
}
