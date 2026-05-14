package natsadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
)

// ErrSubscriberClosed indicates an attempt to Subscribe after Drain.
var ErrSubscriberClosed = errors.New("natsadapter: subscriber closed")

// Subscriber is the minimal subscription surface the adapter consumes from
// the NATS client. Defined here (consumer side) per go-coding-guide.md.
//
// A Subscriber is goroutine-safe: Subscribe may be called concurrently from
// multiple goroutines. Drain blocks until in-flight handlers complete; after
// Drain returns, Subscribe returns ErrSubscriberClosed.
type Subscriber interface {
	// Subscribe registers handler for messages on subject. The returned
	// Unsubscriber stops delivery for this subscription only; Drain stops
	// all subscriptions on the connection.
	//
	// handler runs on a goroutine owned by the Subscriber implementation.
	// The implementation guarantees handler returns before Drain returns.
	Subscribe(ctx context.Context, subject string, handler func(payload []byte)) (Unsubscriber, error)

	// Drain stops all subscriptions, waits for in-flight handlers to
	// return, and closes the underlying connection. Idempotent.
	Drain(ctx context.Context) error
}

// Unsubscriber stops delivery for a single subscription.
type Unsubscriber interface {
	Unsubscribe() error
}

// natsClient is a Subscriber backed by a real NATS connection.
type natsClient struct {
	conn *nats.Conn
}

// Connect dials a NATS broker and returns a Subscriber backed by it.
//
// The caller owns the lifecycle: call Drain to release resources before
// returning from main. Connect respects ctx for the dial timeout; once
// connected the connection's lifetime is tied to Drain, not ctx.
func Connect(ctx context.Context, url string, opts ...nats.Option) (Subscriber, error) {
	conn, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("natsadapter: connect %q: %w", url, err)
	}
	if ctx.Err() != nil {
		// Caller already cancelled; clean up and bail.
		conn.Close()
		return nil, ctx.Err()
	}
	return &natsClient{conn: conn}, nil
}

// Subscribe registers handler on subject. The NATS client invokes handler on
// a per-subscription goroutine; multiple messages on the same subject are
// serialised by the NATS library.
func (c *natsClient) Subscribe(ctx context.Context, subject string, handler func(payload []byte)) (Unsubscriber, error) {
	if c.conn.IsClosed() {
		return nil, ErrSubscriberClosed
	}
	sub, err := c.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	if err != nil {
		return nil, fmt.Errorf("natsadapter: subscribe %q: %w", subject, err)
	}
	return sub, nil
}

// Drain stops all subscriptions and closes the connection. Drain blocks until
// all in-flight handlers return or ctx is cancelled, whichever comes first.
//
// Per nats.go semantics, Drain is preferable to Close because it lets
// in-flight messages finish; Close discards them.
func (c *natsClient) Drain(ctx context.Context) error {
	done := make(chan error, 1)
	go func() {
		done <- c.conn.Drain()
	}()
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("natsadapter: drain: %w", err)
		}
		return nil
	case <-ctx.Done():
		// Best-effort: force close to unblock the caller. The drain
		// goroutine will eventually finish; we just stop waiting.
		c.conn.Close()
		return fmt.Errorf("natsadapter: drain cancelled: %w", ctx.Err())
	}
}
