package monitor

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

func getUptime() (time.Duration, error) {
	var d time.Duration

	uptime, err := os.ReadFile(path.Join("/proc", "uptime"))
	if err != nil {
		return d, err
	}

	seconds := bytes.Fields(uptime)[0]

	s, err := strconv.ParseFloat(string(seconds), 32)
	if err != nil {
		return d, err
	}

	d = time.Duration(s) * time.Second

	return d, nil
}

func FormatUptime(d time.Duration) string {
	h := int(d.Hours()) % 60
	m := int(d.Minutes()) % 60

	return fmt.Sprintf("%02dh%02dm", h, m)
}
