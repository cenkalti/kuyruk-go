package kuyruk

import (
	"errors"
	"time"

	"github.com/cenkalti/redialer/amqpredialer"
	"github.com/streadway/amqp"
)

// Kuyruk is the main type of this library.
// It holds a single connection to RabbitMQ server.
type Kuyruk struct {
	config Config
	amqp   *amqpredialer.AMQPRedialer
}

// New returns a new Kuyruk value, validating the config.
func New(config Config) (*Kuyruk, error) {
	r, err := amqpredialer.New(config.RabbitURI)
	if err != nil {
		return nil, err
	}
	return &Kuyruk{
		config: config,
		amqp:   r,
	}, nil
}

// Run this instance. Call this function with a go statement.
// It will connect to the RabbitMQ server.
func (k *Kuyruk) Run() {
	k.amqp.Run()
}

func (k *Kuyruk) getConnectionWithTimeout() (*amqp.Connection, error) {
	select {
	case conn := <-k.amqp.Conn():
		return conn, nil
	case <-time.After(k.config.RabbitConnectionWaitTimeout):
		return nil, errors.New("timeout")
	}
}
