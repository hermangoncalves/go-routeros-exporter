package prometheus

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hermangoncalves/go-routeros-exporter/ports"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// cpuUsage = prometheus.NewGauge(
	// 	prometheus.GaugeOpts{
	// 		Name: "mikrotik_cpu_usage",
	// 		Help: "CPU usage percentage",
	// 	},
	// )
	// memoryUsage = prometheus.NewGauge(
	// 	prometheus.GaugeOpts{
	// 		Name: "mikrotik_memory_usage",
	// 		Help: "Memory usage percentage",
	// 	},
	// )
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
}

func NewMetricsHandler(metricsService ports.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsService: metricsService}
}

func (h *MetricsHandler) UpdateMetrics() error {
	metrics, err := h.metricsService.CollectMetrics()
	if err != nil {
		return err
	}

	log.Println(metrics.InterfaceTraffic)

	for iface, traffic := range metrics.InterfaceTraffic {
		rxBytes, err := strconv.ParseFloat(strings.TrimSpace(traffic.RxBytes), 64)
		if err != nil {
			log.Printf("failed to convert RxBytes: %v", err)
		}
		txBytes, err := strconv.ParseFloat(strings.TrimSpace(traffic.TxBytes), 64)
		if err != nil {
			log.Printf("failed to convert TxBytes: %v", err)
		}
		interfaceRxBytes.WithLabelValues(iface).Set(float64(rxBytes))
		interfaceTxBytes.WithLabelValues(iface).Set(float64(txBytes))
	}

	CPUUsage, err := strconv.ParseFloat(strings.TrimSpace(metrics.CPUUsage), 64)
	if err != nil {
		log.Printf("failed to convert CPUUsage: %v", err)
	}

	MemoryUsage, err := strconv.ParseFloat(strings.TrimSpace(metrics.MemoryUsage), 64)
	if err != nil {
		log.Printf("failed to convert CPUUsage: %v", err)
	}

	systemMetrics.WithLabelValues("cpu").Set(CPUUsage)
	systemMetrics.WithLabelValues("memory").Set(MemoryUsage)

	return nil
}

var promHandler = promhttp.Handler()

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.UpdateMetrics(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	promHandler.ServeHTTP(w, r)
}
