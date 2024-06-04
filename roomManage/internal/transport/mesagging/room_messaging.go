package messaging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"roomManage/internal/domain/model"
)

type RoomMessaging struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRoomMessaging(rabbitMQUrl string) (*RoomMessaging, error) {
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

	return &RoomMessaging{connection: conn, channel: ch}, nil
}

func (m *RoomMessaging) PublishRoomCreated(room *model.Room) error {
	err := m.channel.ExchangeDeclare(
		"room_exchange",
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

	body, err := json.Marshal(room)
	if err != nil {
		return err
	}

	err = m.channel.Publish(
		"room_exchange",
		"room.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish room created message: %v", err)
		return err
	}

	log.Printf("Room created message published: %v", room)
	return nil
}

func (m *RoomMessaging) Close() {
	m.channel.Close()
	m.connection.Close()
}
