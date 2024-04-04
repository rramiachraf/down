package monitor

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"
)

type CPU struct {
	ModelName string `json:"model"`
	Cores     int    `json:"cores"`
}

func getCPUInformation() (CPU, error) {
	var cpu CPU

	f, err := os.Open(path.Join("/proc", "cpuinfo"))
	if err != nil {
		return cpu, err
	}

	defer f.Close()

	sc := bufio.NewScanner(f)
	var c int
	for sc.Scan() {
		c++
		if sc.Err() != nil {
			return cpu, sc.Err()
		}

		key, val, found := strings.Cut(sc.Text(), ":")
		if !found {
			break
		}

		key = cleanValue(key)
		val = cleanValue(val)

		switch key {
		case "model name":
			cpu.ModelName = val
		case "cpu cores":
			if n, err := strconv.Atoi(val); err == nil {
				cpu.Cores = n
			}
		}
	}

	return cpu, nil
}
