package monitor

import (
	"time"
)

type Monitor struct {
	Processes []Process     `json:"processes"`
	Memory    Memory        `json:"memory"`
	Uptime    time.Duration `json:"uptime"`
}

func NewMonitor() *Monitor {
	var m Monitor
	return &m
}

func (m *Monitor) Update() error {
	p, err := getProcesses()
	if err != nil {
		return err
	}

	m.Processes = p

	var mem Memory
	if err := mem.parseData(); err != nil {
		return err
	}

	m.Memory = mem

	uptime, err := getUptime()
	if err != nil {
		return err
	}

	m.Uptime = uptime

	return nil
}
