package monitor

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type Drive struct {
	Identifier string `json:"identifier"`
	Total      int    `json:"total"`
	Used       int    `json:"used"`
}

func getDrives() ([]Drive, error) {
	dp, err := filepath.Glob("/sys/block/sd*")
	if err != nil {
		return nil, err
	}

	var drives []Drive
	for _, f := range dp {
		p := path.Join(f, "size")
		if size, err := os.ReadFile(p); err == nil {
			c := cleanValue(string(size))
			n, err := strconv.Atoi(c)
			if err != nil {
				continue
			}

			total := (n * 512) / 1024

			drives = append(drives, Drive{
				Identifier: path.Base(f),
				Total:      total,
			})
		}
	}

	return drives, nil
}
