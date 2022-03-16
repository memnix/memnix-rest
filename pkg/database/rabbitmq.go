package database

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMqConn *amqp.Connection
var RabbitMqChan *amqp.Channel

func Rabbit() (*amqp.Channel, error) {
	conn, err := amqp.Dial(rabbitMQ)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	fmt.Println("Connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
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
		return nil, fmt.Errorf("failed to declare an exchange: %s", err)
	}

	RabbitMqChan, RabbitMqConn = ch, conn

	return ch, nil
}
func SendMessageToChannel(ch *amqp.Channel, body []byte, key string) error {

	err := ch.Publish(
		"logs", // exchange
		key,    // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to declare send a message: %s", err)
	}
	return nil
}
