package service

import (
	"github.com/hermangoncalves/go-routeros-exporter/core/domain"
	"github.com/hermangoncalves/go-routeros-exporter/ports"
)

type metricsService struct {
	mikrotikClient ports.MikrotikClient
}

func NewNetricsService(mikrotikClient ports.MikrotikClient) ports.MetricsService {
	return &metricsService{mikrotikClient: mikrotikClient}
}

func (s *metricsService) CollectMetrics() (*domain.Metrics, error) {
	interfaceTraffic, err := s.mikrotikClient.GetInterfaceTraffic()
	if err != nil {
		return nil, err
	}

	cpuUsage, memoryUsage, err := s.mikrotikClient.GetSystemResources()
	if err != nil {
		return nil, err
	}

	return &domain.Metrics{
		InterfaceTraffic: interfaceTraffic,
		CPUUsage:         cpuUsage,
		MemoryUsage:      memoryUsage,
	}, nil
}
