package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/go-routeros/routeros/v3"
	"github.com/hermangoncalves/go-routeros-exporter/core/domain"
	"github.com/hermangoncalves/go-routeros-exporter/ports"
)

type MikrotikClient struct {
	client *routeros.Client
}

func NewMikrotikClient(ctx context.Context, host, username, password string) (ports.MikrotikClient, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		client *routeros.Client
		err    error
	}
	ch := make(chan result, 1)

	go func() {
		client, err := routeros.DialContext(ctx, host, username, password)

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

		return &MikrotikClient{res.client}, nil

	case <-ctx.Done():
		// Handle context cancellation or timeout
		return nil, fmt.Errorf("%s: %v", "Network Timeout", ctx.Err())
	}
}

func (c *MikrotikClient) GetInterfaceTraffic() (map[string]domain.InterfaceTraffic, error) {
	reply, err := c.client.Run("/interface/monitor-traffic", "=interface=ether1", "=once=")
	if err != nil {
		return nil, fmt.Errorf("failed to get interface traffic: %w", err)
	}

	traffic := make(map[string]domain.InterfaceTraffic)

	for _, re := range reply.Re {
		name := re.Map["name"]
		rxBytes := re.Map["rx-bits-per-second"]
		txBytes := re.Map["tx-bits-per-second"]
		traffic[name] = domain.InterfaceTraffic{
			RxBytes: rxBytes,
			TxBytes: txBytes,
		}
	}
	return traffic, nil
}

func (c *MikrotikClient) GetSystemResources() (cpuUsage string, memoryUsage string, err error) {
	reply, err := c.client.Run("/system/resource/print")
	if err != nil {
		return "", "", fmt.Errorf("failed to get system resources: %w", err)
	}

	cpuUsage = reply.Re[0].Map["cpu-load"]
	memoryUsage = reply.Re[0].Map["free-memory"]
	return cpuUsage, memoryUsage, nil
}
