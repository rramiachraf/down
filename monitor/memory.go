package monitor

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"
)

type Memory struct {
	Total        int `json:"total"`
	Free         int `json:"free"`
	Available    int `json:"available"`
	Buffers      int `json:"buffers"`
	Cached       int `json:"cached"`
	Shared       int `json:"shared"`
	SReclaimable int `json:"sreclaimable"`
	SwapTotal    int `json:"swapTotal"`
	SwapFree     int `json:"swapFree"`
}

func (m *Memory) parseData() error {
	f, err := os.Open(path.Join("/proc", "meminfo"))
	if err != nil {
		return err
	}

	defer f.Close()

	b := bufio.NewScanner(f)
	for b.Scan() {
		if b.Err() != nil {
			return err
		}

		entry := b.Text()
		key, value, found := strings.Cut(entry, ":")
		if !found {
			continue
		}

		value = cleanValue(value)

		switch key {
		case "MemTotal":
			m.Total = parseMemoryValue(value)
		case "MemFree":
			m.Free = parseMemoryValue(value)
		case "MemAvailable":
			m.Available = parseMemoryValue(value)
		case "Buffers":
			m.Buffers = parseMemoryValue(value)
		case "Cached":
			m.Cached = parseMemoryValue(value)
		case "Shmem":
			m.Shared = parseMemoryValue(value)
		case "SReclaimable":
			m.SReclaimable = parseMemoryValue(value)
		case "SwapTotal":
			m.SwapTotal = parseMemoryValue(value)
		case "SwapFree":
			m.SwapFree = parseMemoryValue(value)
		}
	}

	return nil
}

func (m Memory) CalculateUsed() int {
	return (m.Total + m.Shared - m.Free - m.Buffers - m.Cached - m.SReclaimable)
}

func parseMemoryValue(v string) int {
	n, _ := strconv.Atoi(v[:len(v)-3])
	return n
}

func ConvertToGib(kb int) float32 {
	return float32(kb) / 1024 / 1024
}

func cleanValue(v string) string {
	clean := strings.ReplaceAll(v, "\x00", " ")

	return strings.TrimSpace(clean)
}
