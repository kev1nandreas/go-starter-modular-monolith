package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kev1nandreas/go-rest-api-template/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ holds a connection and a channel for publishing/consuming.
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	host := config.GetEnvString("RABBITMQ_HOST", "localhost")
	port := config.GetEnvString("RABBITMQ_PORT", "5672")
	user := config.GetEnvString("RABBITMQ_USER", "guest")
	password := config.GetEnvString("RABBITMQ_PASSWORD", "guest")
	vhost := config.GetEnvString("RABBITMQ_VHOST", "/")

	u := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   "/",
	}
	if vhost != "/" {
		u.Path = "/" + vhost
	}
	if config.GetEnvBool("RABBITMQ_USE_SSL", false) {
		u.Scheme = "amqps"
	}

	conn, err := amqp.Dial(u.String())
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}
	return &RabbitMQ{Conn: conn, Channel: ch}, nil
}

// DeclareQueue declares a durable queue.
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.Channel.QueueDeclare(name, true, false, false, false, nil)
}

// Publish sends a persistent message to a queue via the default exchange.
func (r *RabbitMQ) Publish(ctx context.Context, queue string, body []byte, contentType string) error {
	return r.Channel.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:  contentType,
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

// PublishJSON marshals v to JSON and publishes it to a queue.
func (r *RabbitMQ) PublishJSON(ctx context.Context, queue string, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("rabbitmq marshal: %w", err)
	}
	return r.Publish(ctx, queue, body, "application/json")
}

// Consume starts consuming a queue with manual ack.
func (r *RabbitMQ) Consume(queue string) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(queue, "", false, false, false, false, nil)
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}
