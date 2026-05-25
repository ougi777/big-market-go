package rabbitmq

import (
	"context"
	"errors"
	"testing"
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
