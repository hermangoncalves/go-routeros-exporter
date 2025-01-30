package ports

import (
	"context"

	"github.com/go-routeros/routeros/v3"
)

type MikrotikClient interface {
	RunCommand(ctx context.Context, command string, args ...string) (*routeros.Reply, error)
	GetSystemIdentity() (string, error)
	Close() error
}

type RouterosClient interface {
	RunArgs(args []string) (*routeros.Reply, error)
	Close() error
}

type MikrotikAuthenticator interface {
	Authenticate(ctx context.Context, address, username, password string) (MikrotikClient, error)
}
