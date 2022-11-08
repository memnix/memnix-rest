package infrastructures

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMq struct {
	Channel     *amqp.Channel
	Connection  *amqp.Connection
	rabbitMqUrl string
}

var (
	// RabbitMQ is the rabbitmq connection handler
	RabbitMQ RabbitMq
)

func Rabbit() (*amqp.Channel, error) {

	RabbitMQ = loadRabbitMQVar()

	conn, err := amqp.Dial(RabbitMQ.rabbitMqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	log.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		"logs",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	RabbitMQ.Channel, RabbitMQ.Connection = ch, conn

	return ch, nil
}

func SendMessageToChannel(ch *amqp.Channel, body []byte, key string) error {
	err := ch.PublishWithContext(
		context.Background(),
		"logs", // exchange
		key,    // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to declare send a message: %w", err)
	}
	return nil
}
