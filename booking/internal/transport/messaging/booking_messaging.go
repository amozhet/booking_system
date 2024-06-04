package messaging

import (
	"encoding/json"
	"log"

	"booking/internal/domain/model"
	"github.com/streadway/amqp"
)

type BookingMessaging interface {
	PublishBookingCreated(booking *model.Booking) error
}
type BookingMessagingImpl struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewBookingMessaging(rabbitMQUrl string) (BookingMessaging, error) {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return nil, err
	}

	return &BookingMessagingImpl{connection: conn, channel: ch}, nil
}

func (m *BookingMessagingImpl) PublishBookingCreated(booking *model.Booking) error {
	err := m.channel.ExchangeDeclare(
		"booking_exchange", // exchange name
		"topic",            // exchange type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
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

func (m *BookingMessagingImpl) Close() {
	m.channel.Close()
	m.connection.Close()
}
