package messages

import "time"

type AboutNode struct {
	ID          string
	Name        string
	Uptime      time.Duration
	LastContact time.Duration
	CPUUsage    float64
	MemoryUsage float64
	Leader      bool
	Version     string
}

type AboutMsg struct {
	Nodes       []AboutNode
	ID          string
	Name        string
	Version     string
	Degraded    bool
	DegradedErr string
}

type ErrorMsg error
