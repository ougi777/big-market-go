package rabbitmq

import (
	"context"

	"bm-go/internal/config"
	"bm-go/internal/domain/award"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection
}

var _ award.MessagePublisher = (*Client)(nil)

func Dial(cfg config.RabbitMQConfig) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) Publish(ctx context.Context, topic string, message string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	if _, err := ch.QueueDeclare(topic, true, false, false, false, nil); err != nil {
		return err
	}
	return ch.PublishWithContext(ctx, "", topic, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(message),
	})
}
