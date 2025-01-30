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
	Client ports.RouterosClient
}

type MikrotikAuthAdapter struct {
	timeout time.Duration
}

func NewMikrotikAuthenticator(timeout time.Duration) ports.MikrotikAuthenticator {
	return &MikrotikAuthAdapter{
		timeout: timeout,
	}
}

func (m *MikrotikAuthAdapter) Authenticate(ctx context.Context, address, username, password string) (client ports.MikrotikClient, err error) {
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

		// Avoid sending if the context is already canceled
		select {
		case ch <- result{client, err}:
		case <-ctx.Done(): // Exit goroutine if context is canceled
			if client != nil {
				client.Close() // Cleanup to prevent resource leak
			}
		}
	}()

	// wait for either the connection attempt or the context to timeout
	select {
	case res := <-ch:
		if res.err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidCredentials, res.err)
		}
		if res.client == nil {
			return nil, fmt.Errorf("%w: client is nil", ErrDeviceUnreachable)
		}

		return &MikrotikClientAdapter{res.client}, nil

	case <-ctx.Done():
		// Handle context cancellation or timeout
		return nil, fmt.Errorf("%w: %v", ErrNetworkTimeout, ctx.Err())
	}
}

func (m *MikrotikClientAdapter) RunCommand(ctx context.Context, command string, args ...string) (*routeros.Reply, error) {
	return m.Client.RunArgs(append([]string{command}, args...))
}

func (m *MikrotikClientAdapter) GetSystemIdentity() (string, error) {
	reply, err := m.RunCommand(context.Background(), "/system/identity/print")
	if err != nil {
		return "", fmt.Errorf("failed to get system identity: %w", err)
	}

	if len(reply.Re) == 0 {
		return "", errors.New("failed to get system identity")
	}

	identify, exists := reply.Re[0].Map["name"]
	if !exists {
		return "", fmt.Errorf("missing 'name' field in response")
	}

	return identify, nil
}

func (m *MikrotikClientAdapter) Close() error {
	return m.Client.Close()
}
