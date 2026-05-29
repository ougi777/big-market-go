package rabbitmq

import (
	"context"
	"errors"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestClientPublishWithoutConnection(t *testing.T) {
	var client *Client

	err := client.Publish(context.Background(), "topic", "{}")
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}

func TestClientConsumeWithoutConnection(t *testing.T) {
	var client *Client

	err := client.Consume(context.Background(), "topic", func(context.Context, string) error { return nil })
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}

func TestClientCloseWithoutConnection(t *testing.T) {
	var client *Client

	if err := client.Close(); err != nil {
		t.Fatalf("expected nil close error, got %v", err)
	}

	client = &Client{}
	if err := client.Close(); err != nil {
		t.Fatalf("expected nil close error, got %v", err)
	}
}

func TestClientPublishWithConfirmAck(t *testing.T) {
	channel := &fakePublishChannel{ack: true}
	client := &Client{
		publishChannel:        func() (publishChannel, error) { return channel, nil },
		publishRetryDelays:    []time.Duration{},
		publishConfirmTimeout: time.Second,
	}

	err := client.Publish(context.Background(), "topic", "{}")
	if err != nil {
		t.Fatalf("expected publish success, got %v", err)
	}
	if !channel.confirmEnabled {
		t.Fatal("expected confirm enabled")
	}
	if channel.publishing.DeliveryMode != amqp.Persistent {
		t.Fatalf("expected persistent delivery mode, got %d", channel.publishing.DeliveryMode)
	}
	if channel.publishing.ContentType != "application/json" {
		t.Fatalf("expected json content type, got %s", channel.publishing.ContentType)
	}
	if channel.key != "topic" || string(channel.publishing.Body) != "{}" {
		t.Fatalf("expected topic message, got %s/%s", channel.key, string(channel.publishing.Body))
	}
}

func TestClientPublishRetriesAfterNack(t *testing.T) {
	attempts := 0
	client := &Client{
		publishRetryDelays:    []time.Duration{0},
		publishConfirmTimeout: time.Second,
	}
	client.publishChannel = func() (publishChannel, error) {
		attempts++
		return &fakePublishChannel{ack: attempts == 2}, nil
	}

	err := client.Publish(context.Background(), "topic", "{}")
	if err != nil {
		t.Fatalf("expected publish success after retry, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestClientPublishReturnsNackAfterRetries(t *testing.T) {
	attempts := 0
	client := &Client{
		publishRetryDelays:    []time.Duration{0, 0},
		publishConfirmTimeout: time.Second,
	}
	client.publishChannel = func() (publishChannel, error) {
		attempts++
		return &fakePublishChannel{ack: false}, nil
	}

	err := client.Publish(context.Background(), "topic", "{}")
	if !errors.Is(err, ErrPublishNacked) {
		t.Fatalf("expected nack error, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestClientPublishConfirmTimeout(t *testing.T) {
	client := &Client{
		publishChannel: func() (publishChannel, error) {
			return &fakePublishChannel{skipConfirm: true}, nil
		},
		publishRetryDelays:    []time.Duration{},
		publishConfirmTimeout: time.Millisecond,
	}

	err := client.Publish(context.Background(), "topic", "{}")
	if !errors.Is(err, ErrPublishConfirmTimeout) {
		t.Fatalf("expected confirm timeout, got %v", err)
	}
}

type fakePublishChannel struct {
	ack            bool
	skipConfirm    bool
	confirmEnabled bool
	confirmCh      chan amqp.Confirmation
	key            string
	publishing     amqp.Publishing
}

func (f *fakePublishChannel) Close() error {
	return nil
}

func (f *fakePublishChannel) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{Name: name}, nil
}

func (f *fakePublishChannel) Confirm(noWait bool) error {
	f.confirmEnabled = true
	return nil
}

func (f *fakePublishChannel) NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation {
	f.confirmCh = confirm
	return confirm
}

func (f *fakePublishChannel) PublishWithContext(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	f.key = key
	f.publishing = msg
	if f.skipConfirm {
		return nil
	}
	f.confirmCh <- amqp.Confirmation{Ack: f.ack}
	return nil
}
