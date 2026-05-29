package rabbitmq

import (
	"context"
	"errors"
	"time"

	"bm-go/internal/config"
	"bm-go/internal/domain/award"

	amqp "github.com/rabbitmq/amqp091-go"
)

const defaultPublishConfirmTimeout = 5 * time.Second

var defaultPublishRetryDelays = []time.Duration{
	50 * time.Millisecond,
	100 * time.Millisecond,
	200 * time.Millisecond,
}

type Client struct {
	conn                  *amqp.Connection
	publishChannel        func() (publishChannel, error)
	publishRetryDelays    []time.Duration
	publishConfirmTimeout time.Duration
}

type publishChannel interface {
	Close() error
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	Confirm(noWait bool) error
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation
	PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
}

var ErrClientNotConnected = errors.New("rabbitmq client is not connected")
var ErrPublishNacked = errors.New("rabbitmq publish returned nack")
var ErrPublishConfirmTimeout = errors.New("rabbitmq publish confirm timeout")

var _ award.MessagePublisher = (*Client)(nil)

func Dial(cfg config.RabbitMQConfig) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:                  conn,
		publishRetryDelays:    defaultPublishRetryDelays,
		publishConfirmTimeout: defaultPublishConfirmTimeout,
	}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) Publish(ctx context.Context, topic string, message string) error {
	if c == nil || (c.conn == nil && c.publishChannel == nil) {
		return ErrClientNotConnected
	}
	var err error
	for attempt := 0; ; attempt++ {
		err = c.publishOnce(ctx, topic, message)
		if err == nil {
			return nil
		}

		if isContextError(err) || attempt >= len(c.retryDelays()) {
			return err
		}

		//退让重试
		if sleepErr := sleepWithContext(ctx, c.retryDelays()[attempt]); sleepErr != nil {
			return sleepErr
		}
	}
}

func (c *Client) publishOnce(ctx context.Context, topic string, message string) error {
	ch, err := c.openPublishChannel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	//声明一个队列quene,第二个参数意思是队列持久化
	if _, err := ch.QueueDeclare(topic, true, false, false, false, nil); err != nil {
		return err
	}

	//开启生产者confirm机制
	if err := ch.Confirm(false); err != nil {
		return err
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	if err := ch.PublishWithContext(ctx, "", topic, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(message),
	}); err != nil {
		return err
	}

	timer := time.NewTimer(c.confirmTimeout())
	defer timer.Stop()

	select {
	case confirm := <-confirms:
		if confirm.Ack {
			return nil
		}
		return ErrPublishNacked
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return ErrPublishConfirmTimeout
	}
}

func (c *Client) openPublishChannel() (publishChannel, error) {
	if c.publishChannel != nil {
		return c.publishChannel()
	}
	return c.conn.Channel()
}

func (c *Client) retryDelays() []time.Duration {
	if c == nil || c.publishRetryDelays == nil {
		return defaultPublishRetryDelays
	}
	return c.publishRetryDelays
}

func (c *Client) confirmTimeout() time.Duration {
	if c == nil || c.publishConfirmTimeout <= 0 {
		return defaultPublishConfirmTimeout
	}
	return c.publishConfirmTimeout
}

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func isContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func (c *Client) Consume(ctx context.Context, topic string, handler func(context.Context, string) error) error {
	if c == nil || c.conn == nil {
		return ErrClientNotConnected
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(topic, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return err
	}
	//第三个参数意思是autoack关闭，需要业务方收到那个ack
	deliveries, err := ch.Consume(topic, "", false, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		return err
	}

	go func() {
		defer func() { _ = ch.Close() }()
		for {
			select {
			case <-ctx.Done():
				return
			case delivery, ok := <-deliveries:
				if !ok {
					return
				}
				if err := handler(ctx, string(delivery.Body)); err != nil {
					_ = delivery.Nack(false, true) //第二个参数意思是requeue=true重新入队
					continue
				}
				_ = delivery.Ack(false)
			}
		}
	}()
	return nil
}
