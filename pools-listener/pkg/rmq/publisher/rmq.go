package publisher

import (
	"encoding/json"
	"fmt"
	"log"

	"main/pkg/rmq/event"

	"github.com/streadway/amqp"
)

// RabbitMQSender представляет отправителя сообщений в RabbitMQ
type RabbitMQSender struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQSender создает нового отправителя сообщений в RabbitMQ
func NewRabbitMQSender(amqpURL, queueName string) (*RabbitMQSender, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &RabbitMQSender{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

// SendMessage отправляет сообщение в очередь RabbitMQ
func (r *RabbitMQSender) SendMessage(message string) error {
	err := r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Message sent: %s", message)
	return nil
}

// Close закрывает соединение и канал
func (r *RabbitMQSender) Close() {
	r.channel.Close()
	r.conn.Close()
}

const (
	amqpURL   = "amqp://guest:guest@localhost:5672/"
	queueName = "mm_events"
)

func SendEvent(data interface{}) {
	sender, err := NewRabbitMQSender(amqpURL, queueName)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ sender: %s", err)
	}
	defer sender.Close()

	message := event.BaseMessage{
		Event: "new_price",
		Data:  data,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %s", err)
	}

	err = sender.SendMessage(string(messageBytes))
	if err != nil {
		log.Fatalf("Failed to send message: %s", err)
	}
}
