package monitor

import (
	"strings"
	"time"
)

type Monitor struct {
	Processes []Process     `json:"processes"`
	Memory    Memory        `json:"memory"`
	Drives    []Drive       `json:"drives"`
	CPU       CPU           `json:"cpu"`
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

	d, err := getDrives()
	if err != nil {
		return err
	}

	m.Drives = d

	cpu, err := getCPUInformation()
	if err != nil {
		return err
	}

	m.CPU = cpu

	return nil
}

func ConvertToGib(kb int) float32 {
	return float32(kb) / 1024 / 1024
}

func cleanValue(v string) string {
	clean := strings.ReplaceAll(v, "\x00", " ")

	return strings.TrimSpace(clean)
}
