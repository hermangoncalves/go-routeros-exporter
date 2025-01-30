package ports

import (
	"context"

	"github.com/go-routeros/routeros/v3"
)

type MikrotikClient interface {
	RunCommand(command string, args ...string) (*routeros.Reply, error)
	Close() error
}

type AuthenticationPort interface {
	Authenticate(ctx context.Context, username, password string) (MikrotikClient, error)
}
