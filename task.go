package kuyruk

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

var localhost string

func init() {
	var err error
	localhost, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
}

// Task that is going to be run on a Kuyruk worker.
type Task struct {
	kuyruk   *Kuyruk
	module   string
	function string
	queue    string
}

// Task creates a definition to run a function on a remote worker.
func (k *Kuyruk) Task(module, function, queue string) *Task {
	return &Task{
		kuyruk:   k,
		module:   module,
		function: function,
		queue:    queue,
	}
}

// SendToQueue serializes the task and publishes a message to RabbitMQ server.
func (t *Task) SendToQueue(args []interface{}, kwargs map[string]interface{}) error {
	return t.SendToQueueHost(args, kwargs, "")
}

// SendToQueueHost works like SendToQueue but allows appending hostname to the queue name.
func (t *Task) SendToQueueHost(args []interface{}, kwargs map[string]interface{}, hostname string) error {
	if args == nil {
		args = make([]interface{}, 0)
	}
	if kwargs == nil {
		kwargs = make(map[string]interface{})
	}

	con, err := t.kuyruk.getConnectionWithTimeout()
	if err != nil {
		return err
	}
	ch, err := con.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue := t.queue
	if hostname == "localhost" {
		hostname = localhost
	}
	if hostname != "" {
		queue += "." + hostname
	}
	_, err = ch.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}
	task := map[string]interface{}{
		"id":               uuid.NewV1().String(),
		"module":           t.module,
		"function":         t.function,
		"args":             args,
		"kwargs":           kwargs,
		"queue":            queue,
		"sender_hostname":  localhost,
		"sender_pid":       os.Getpid(),
		"sender_cmd":       strings.Join(os.Args, " "),
		"sender_timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}
	p := amqp.Publishing{
		Body: body,
	}
	return ch.Publish(
		"",    // exchange
		queue, // key
		true,  // mandatory
		false, // immedeate
		p,
	)
}
