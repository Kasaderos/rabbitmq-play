package main

import (

	"github.com/streadway/amqp"
)

const (
	ConnStr       = "amqp://user:password@127.0.0.1:5672/"
	PrefetchCount = 5
	Exchange      = "exch"
	DeadExchange  = "dead-exch"
	QueueName     = "queue"
	DeadQueueName = "dead-queue"
	RoutingKey    = "routingkey"
)

// RMQReceiver receives protocol messages from RabbitMQ
type RMQReceiver struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	DeliveryChan <-chan amqp.Delivery
	NotifyChan   <-chan *amqp.Error
}

// NewRMQReceiver constructs new RabbitMQ receiver
func NewRMQReceiver() (*RMQReceiver, error) {
	r := &RMQReceiver{}
	return r, r.connect()
}

// connect creates connection to RabbitMQ, opens a channel and declares direct exchange
func (r *RMQReceiver) connect() error {
	// create a connection to rabbit
	conn, err := amqp.Dial(ConnStr)
	if err != nil {
		return err
	}

	// create a channel from a connection
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	err = channel.Qos(PrefetchCount, 0, false)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	// declare exchanges
	err = channel.ExchangeDeclare(Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}
	err = channel.ExchangeDeclare(DeadExchange, "direct", true, false, false, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	// declare queues
	_, err = channel.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-queue-mode":              "lazy",
			"x-dead-letter-exchange":    DeadExchange,
			"x-dead-letter-routing-key": RoutingKey,
		},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}
	_, err = channel.QueueDeclare(
		DeadQueueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange": Exchange, // our default exchange
		},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	// Bindings
	err = channel.QueueBind(QueueName, RoutingKey, Exchange, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}
	// TODO: should the dead queque be used?
	err = channel.QueueBind(DeadQueueName, RoutingKey, DeadExchange, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	deliveryChan, err := channel.Consume(QueueName, "", false, false, false, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	r.NotifyChan = conn.NotifyClose(make(chan *amqp.Error))
	r.conn = conn
	r.channel = channel
	r.DeliveryChan = deliveryChan
	return nil
}

// Close closes the resources at its disposal
func (r *RMQReceiver) Close() error {
	r.channel.Close()
	return r.conn.Close()
}
