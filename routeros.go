package main

import (
	"context"
	"fmt"
	"log"
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

func (r *RouterosClient) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}

func fetchMetrics(client *routeros.Client) {
	resp, err := client.Run("/interface/print", "=.proplist=name,tx-byte")
	if err != nil {
		log.Printf("Failed to fetch interface metrics: %v", err)
		return
	}

	for _, re := range resp.Re {
		name := re.Map["name"]
		txBytes.WithLabelValues(name).Set(parseFloat(re.Map["tx-byte"]))
	}
}

func parseFloat(value string) float64 {
	var result float64
	fmt.Sscanf(value, "%f", &result)
	return result
}

func startMetricsCollection(client *routeros.Client, interval time.Duration) {
	for {
		fetchMetrics(client)
		time.Sleep(interval)
	}
}
