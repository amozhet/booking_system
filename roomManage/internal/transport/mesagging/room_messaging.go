package messaging

import (
	"encoding/json"
	"log"
	"roomManage/internal/config"
	"roomManage/internal/domain/model"
	"roomManage/internal/service"

	"github.com/streadway/amqp"
)

type RoomMessaging struct {
	service *service.RoomService
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRoomMessaging(cfg *config.Config, service *service.RoomService) (*RoomMessaging, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RoomMessaging{
		service: service,
		conn:    conn,
		channel: ch,
	}, nil
}

func (m *RoomMessaging) ConsumeMessages() error {
	msgs, err := m.channel.Consume(
		"room_queue", // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var room model.Room
		if err := json.Unmarshal(msg.Body, &room); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// Process the room message
		if err := m.processRoomMessage(&room); err != nil {
			log.Printf("Error processing room message: %v", err)
		}
	}

	return nil
}

func (m *RoomMessaging) processRoomMessage(room *model.Room) error {
	// Implement the logic to process the room message
	// This could include saving the room to the database, updating existing room, etc.
	// Here is a simple example of saving the room to the database:

	if err := m.service.CreateRoom(room); err != nil {
		return err
	}
	return nil
}

func (m *RoomMessaging) PublishRoomMessage(room *model.Room) error {
	body, err := json.Marshal(room)
	if err != nil {
		return err
	}

	err = m.channel.Publish(
		"",           // exchange
		"room_queue", // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}

func (m *RoomMessaging) Close() {
	m.channel.Close()
	m.conn.Close()
}
