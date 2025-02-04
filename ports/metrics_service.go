package ports

import "github.com/hermangoncalves/go-routeros-exporter/core/domain"

type MetricsService interface {
	CollectMetrics() (*domain.Metrics, error)
}