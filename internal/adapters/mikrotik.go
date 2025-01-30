package adapters

import (
	"context"

	"github.com/go-routeros/routeros/v3"
	"github.com/hermangoncalves/go-routeros-exporter/internal/ports"
)

type MikrotikClientAdapter struct {
	Client *routeros.Client
}

type MikrotikAuthAdapter struct {
	address string
}

func NewMikrotikAuthAdapter(address string) *MikrotikAuthAdapter {
	return &MikrotikAuthAdapter{
		address: address,
	}
}

func (m *MikrotikAuthAdapter) Authenticate(ctx context.Context, username, password string) (ports.MikrotikClient, error) {
	client, err := routeros.DialContext(ctx, m.address, username, password)
	if err != nil {
		return nil, err
	}

	return &MikrotikClientAdapter{Client: client}, nil
}

func (m *MikrotikClientAdapter) RunCommand(command string, args ...string) (*routeros.Reply, error) {
	return m.Client.RunArgs(append([]string{command}, args...))
}

func (m *MikrotikClientAdapter) Close() error {
	return m.Client.Close()
}
