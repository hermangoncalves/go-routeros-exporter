package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-routeros/routeros/v3"
)

type RouterosClient struct {
	Client *routeros.Client
}

func NewRouterosClient(address, username, password string) (*RouterosClient, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		client *routeros.Client
		err    error
	}
	ch := make(chan result, 1)

	go func() {
		client, err := routeros.DialContext(ctx, address, username, password)

		select {
		case ch <- result{client, err}:
		case <-ctx.Done():
			if client != nil {
				client.Close()
			}
		}
	}()

	select {
	case res := <-ch:
		if res.err != nil {
			return nil, fmt.Errorf("%s: %w", "Invalid Credentials", res.err)
		}
		if res.client == nil {
			return nil, fmt.Errorf("%s: client is nil", "Device Unreachable")
		}

		return &RouterosClient{res.client}, nil

	case <-ctx.Done():
		// Handle context cancellation or timeout
		return nil, fmt.Errorf("%s: %v", "Network Timeout", ctx.Err())
	}
}
