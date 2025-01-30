package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-routeros/routeros/v3"
	"github.com/hermangoncalves/go-routeros-exporter/internal/ports"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDeviceUnreachable  = errors.New("device unreachable")
	ErrNetworkTimeout     = errors.New("network timeout")
)

type MikrotikClientAdapter struct {
	Client *routeros.Client
}

type MikrotikAuthAdapter struct {
	timeout time.Duration
}

func NewMikrotikAuthAdapter(timeout time.Duration) *MikrotikAuthAdapter {
	return &MikrotikAuthAdapter{
		timeout: timeout,
	}
}

func (m *MikrotikAuthAdapter) Authenticate(ctx context.Context, address, username, password string) (ports.MikrotikClient, error) {
	// Created a new context with the specified timeout
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel() // Ensure the context is canceled to release resources

	// create a channel to handle the connectin attempt
	type result struct {
		client *routeros.Client
		err    error
	}
	ch := make(chan result, 1)

	go func() {
		client, err := routeros.DialContext(ctx, address, username, password)
		ch <- result{client, err}
	}()

	// wait for either the connection attempt or the context to timeout
	select {
	case res := <-ch:
		if res.err != nil {
			return nil, fmt.Errorf("authentication to Mikrotik Device: %w", res.err)
		}

		return &MikrotikClientAdapter{res.client}, nil

	case <-ctx.Done():
		// Handle context cancellation or timeout
		return nil, fmt.Errorf("%w: %v", ErrNetworkTimeout, ctx.Err())
	}
}

func (m *MikrotikClientAdapter) RunCommand(command string, args ...string) (*routeros.Reply, error) {
	return m.Client.RunArgs(append([]string{command}, args...))
}

func (m *MikrotikClientAdapter) Close() error {
	return m.Client.Close()
}
