package client

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/rabbitmq/amqp091-go"
)

// ProcessFunc is a function that processes messages.
type ProcessFunc func(message string) error

// MessageSender is an interface for sending messages.
type MessageSender interface {
	Send(queue string, message string, cxt context.Context) error
}

// MessageReceiver is an interface for receiving messages
type MessageReceiver interface {
	StartReceiving(queue string, ctx context.Context)
}

// RMQClient is a client for sending messages to rmq message queues
type RMQClient struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
	l    *slog.Logger
}

func NewRabbitMQClient(config *config.RMQConfig) (*RMQClient, error) {
	l := config.L.With()
	connString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
	conn, err := amqp091.Dial(connString)
	if err != nil {
		l.Error("Error dialing RMQ")
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		l.Error("Error creating RMQ channel")
		return nil, err
	}
	return &RMQClient{
		conn: conn,
		ch:   ch,
		l:    l,
	}, nil
}

// Send sends a message to the message queue to be consumed by
// receiving applications.
func (rmqc *RMQClient) Send(queue string, message string, ctx context.Context) error {
	_, err := rmqc.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		rmqc.l.Error(
			"Error declaring queue",
			"queue", queue,
			"error", err,
		)

		return err
	}
	err = rmqc.ch.PublishWithContext(ctx, "", queue, false, false, amqp091.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(message),
		DeliveryMode: amqp091.Persistent,
	})
	if err != nil {
		rmqc.l.Error(
			"Error publishing message",
			"queue", queue,
			"error", err,
		)
		return err
	}
	return nil
}

// StartReceiving listens for messages on the queue and processes them.
func (rmqc *RMQClient) StartReceiving(ctx context.Context, queue string, f ProcessFunc) {
	_, err := rmqc.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		rmqc.l.Error("failed to declare queue", "queue", queue, "error", err)
		return
	}
	err = rmqc.ch.Qos(1, 0, false)
	if err != nil {
		rmqc.l.Error("failed to set QoS", "error", err)
		return
	}
	msgs, err := rmqc.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		rmqc.l.Error("failed to start consuming", "queue", queue, "error", err)
		return
	}
	rmqc.l.Info("started receiving messages", "queue", queue)

	// Process messages with context cancellation
	for {
		select {
		case <-ctx.Done():
			rmqc.l.Info("stopping message receiver", "queue", queue, "reason", ctx.Err())
			return
		case msg, ok := <-msgs:
			if !ok {
				rmqc.l.Info("message channel closed", "queue", queue)
				return
			}

			err := f(string(msg.Body))
			if err != nil {
				rmqc.l.Error("failed to process message", "error", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}
}

// Close the client's connections
func (rmqc *RMQClient) Close() error {
	if rmqc.ch != nil {
		rmqc.ch.Close()
	}
	if rmqc.conn != nil {
		return rmqc.conn.Close()
	}
	return nil
}
