package ports

import "github.com/hermangoncalves/go-routeros-exporter/core/domain"

type MikrotikClient interface {
	GetInterfaceTraffic() (map[string]domain.InterfaceTraffic, error)
	GetSystemResources() (cpuUsage string, memoryUsage string, err error)
}
