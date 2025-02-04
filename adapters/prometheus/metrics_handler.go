package prometheus

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hermangoncalves/go-routeros-exporter/ports"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	interfaceRxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mikrotik_interface_rx_bytes",
			Help: "Receive bytes on the interface",
		},
		[]string{"interface"},
	)

	interfaceTxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mikrotik_interface_tx_bytes",
			Help: "Transmitted bytes on the interface",
		},
		[]string{"interface"},
	)

	systemMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mikrotik_system_metrics",
			Help: "Métricas do sistema Mikrotik",
		},
		[]string{"type"},
	)
)

func init() {
	// prometheus.MustRegister(interfaceRxBytes, interfaceTxBytes, cpuUsage, memoryUsage)
	prometheus.MustRegister(interfaceRxBytes, interfaceTxBytes, systemMetrics)
}

type MetricsHandler struct {
	metricsService ports.MetricsService
	logger         *logrus.Logger
}

func NewMetricsHandler(metricsService ports.MetricsService, logger *logrus.Logger) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
		logger:         logger,
	}
}

func (h *MetricsHandler) UpdateMetrics() error {
	metrics, err := h.metricsService.CollectMetrics()
	if err != nil {
		h.logger.WithError(err).Error("Failed to collect metrics")
		return err
	}

	for iface, traffic := range metrics.InterfaceTraffic {
		rxBytes, err := strconv.ParseFloat(strings.TrimSpace(traffic.RxBytes), 64)
		if err != nil {
			h.logger.WithFields(logrus.Fields{
				"interface": iface,
				"rx_bytes":  traffic.RxBytes,
			}).Error("Failed to convert RxBytes")
		}
		txBytes, err := strconv.ParseFloat(strings.TrimSpace(traffic.TxBytes), 64)
		if err != nil {
			h.logger.WithFields(logrus.Fields{
				"interface": iface,
				"tx_bytes":  traffic.TxBytes,
			}).Error("Failed to convert TxBytes")
		}
		interfaceRxBytes.WithLabelValues(iface).Set(float64(rxBytes))
		interfaceTxBytes.WithLabelValues(iface).Set(float64(txBytes))
	}

	CPUUsage, err := strconv.ParseFloat(strings.TrimSpace(metrics.CPUUsage), 64)
	if err != nil {
		h.logger.WithError(err).Error("Failed to convert CPUUsage")
	}

	MemoryUsage, err := strconv.ParseFloat(strings.TrimSpace(metrics.MemoryUsage), 64)
	if err != nil {
		h.logger.WithError(err).Error("Failed to convert MemoryUsage")
	}

	systemMetrics.WithLabelValues("cpu").Set(CPUUsage)
	systemMetrics.WithLabelValues("memory").Set(MemoryUsage)

	return nil
}

var promHandler = promhttp.Handler()

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.UpdateMetrics(); err != nil {
		h.logger.WithError(err).Error("Failed to update metrics before serving HTTP request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	promHandler.ServeHTTP(w, r)
}
