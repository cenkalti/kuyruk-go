package kuyruk

import "time"

// Config for Kuyruk application.
type Config struct {
	RabbitURI                   string
	RabbitConnectionWaitTimeout time.Duration
}

// DefaultConfig is a Config with default values set.
// You should not modify this. Instead copy and modify that.
var DefaultConfig = Config{
	RabbitURI:                   "amqp://guest:guest@localhost:5672/",
	RabbitConnectionWaitTimeout: 30 * time.Second,
}
