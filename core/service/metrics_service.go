package service

import (
	"log"

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

	log.Println(interfaceTraffic["ether1"])

	cpuUsage, memoryUsage, err := s.mikrotikClient.GetSystemResources()
	if err != nil {
		return nil, err
	}

	log.Println("CPU Usage: " + cpuUsage)
	log.Println("Memory Usage: " + memoryUsage)

	return &domain.Metrics{
		InterfaceTraffic: interfaceTraffic,
		CPUUsage:         cpuUsage,
		MemoryUsage:      memoryUsage,
	}, nil
}
