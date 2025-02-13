package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var deviceIndex *int

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel(choices []string) model {
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
	}
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
			deviceIndex = &m.cursor
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var choices []string

	for _, device := range cfg.Devices {
		choices = append(choices, fmt.Sprintf("%s - %s", device.Name, device.Host))
	}

	p := tea.NewProgram(initialModel(choices))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if deviceIndex == nil {
		log.Fatal("No device selected")
	}

	device := cfg.Devices[*deviceIndex]
	log.Println("Connecting to the device...")
	routerosClient, err := NewRouterosClient(
		fmt.Sprintf("%s:%d", device.Host, device.Port),
		device.Username,
		device.Password,
	)
	if err != nil {
		log.Fatalf("failed to create new client: %v", err)
	}

  RunCommand()
  
	startServer(cfg.ListenAddress)

}

func startServer(addr string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Exporter running on http://localhost%s/metrics\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}


func RunCommand() {
reply, err := routerosClient.Run("/system/resources/print")
  if err != nil {
    log.Fatalf("failed to run command: %v", err)
  }
  
fmt.Println(reply.Re)
}
