package main

import (
	"errors"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/streadway/amqp"
)

var (
	// ErrUnsupportedProtocol gets returned if MessageSender cannot send messages of the protocol
	ErrUnsupportedProtocol = errors.New("Unsupported protocol")
	// ErrConnectionClosed gets returned if the connection to queue manager is closed
	ErrConnectionClosed = errors.New("Connection to queue manager is closed")
)

const (
	ExchangeName = "exch"
	Host         = "127.0.0.1"
	Port         = "5672"
	RoutingKey   = "routingkey"
)

// RMQSender is a sender that publishes messages to RabbitMQ
type RMQSender struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	closed  uint32
}

// NewRMQSender constructs new RabbitMQ sender
func NewRMQSender() (*RMQSender, error) {
	s := &RMQSender{}
	return s, s.connect()
}

// Send sends parsing results to RabbitMQ
func (s *RMQSender) Send(msg []byte) error {
	if atomic.LoadUint32(&s.closed) == 1 {
		return ErrConnectionClosed
	}

	err := s.channel.Publish(
		ExchangeName,
		RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers: amqp.Table{
				"x-death-count": 5,
			},
			ContentType:  "application/json",
			Body:         msg,
			DeliveryMode: amqp.Persistent,
		},
	)
	return err
}

// connect creates connection to RabbitMQ, opens a channel and declares direct exchange
func (s *RMQSender) connect() error {
	// create a connection to rabbit
	addr := fmt.Sprintf("amqp://%s:%s@%s:%s/", "user", "password", Host, Port)
	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}

	// create a channel from a connection
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	// create exchange
	err = channel.ExchangeDeclare(ExchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}
	log.Printf("[rmq] instance: %s:%s. connection opened successfully", Host, Port)

	s.conn = conn
	s.channel = channel
	return nil
}

// Close closes the resources at its disposal
func (s *RMQSender) Close() {
	s.channel.Close()
	s.conn.Close()
}
