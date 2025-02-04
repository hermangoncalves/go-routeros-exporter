package domain

type Metrics struct {
	InterfaceTraffic map[string]InterfaceTraffic
	CPUUsage         string
	MemoryUsage      string
}

type InterfaceTraffic struct {
	RxBytes string
	TxBytes string
}
