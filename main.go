package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mikrotikDevices []Device
	selectedDevice  *Device
	txBytes         = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mikrotik_interface_tx_bytes",
			Help: "Transmitted bytes on MikroTik interfaces",
		},
		[]string{"interface"},
	)
)

type model struct {
	choices  []string
	cursor   int
	selected int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.cursor
			selectedDevice = &mikrotikDevices[m.selected]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Select a MikroTik device:\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}

func main() {
	cfg, err := LoadConfig(".config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mikrotikDevices = cfg.Devices

	choices := make([]string, len(cfg.Devices))
	for index, device := range cfg.Devices {
		choices[index] = fmt.Sprintf("%s - %s", device.Name, device.Host)
	}
	p := tea.NewProgram(model{choices: choices})
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running Bubble Tea program: %v", err)
	}

	if selectedDevice == nil {
		log.Fatal("No device selected")
	}

	fmt.Println("Connecting to Mikrotik device...")

	client, err := NewRouterosClient(
		fmt.Sprintf("%s:%d", selectedDevice.Host, selectedDevice.Port),
		selectedDevice.Username,
		selectedDevice.Password,
	)
	if err != nil {
		log.Fatalf("Failed to connect to MikroTik: %v", err)
	}
	defer client.Close()

	go startMetricsCollection(client.Client, 30*time.Second)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("MikroTik Prometheus Exporter running on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))

}

func init() {
	prometheus.MustRegister(
		txBytes,
	)
}
