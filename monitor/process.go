package monitor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
)

type Process struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	Owner       string `json:"owner"`
	Group       string `json:"group"`
	Priority    int    `json:"priority"`
	Nice        int    `json:"nice"`
	Status      byte   `json:"status"`
	VirtualMem  int    `json:"virtualMemory"`
	ResidentMem int    `json:"residentMemory"`
	SharedMem   int    `json:"sharedMemory"`
}

func (p *Process) parseStatus() error {
	if p.ID == "" {
		return fmt.Errorf("pid is undefined")
	}

	f, err := os.Open(path.Join("/proc", p.ID, "status"))
	if err != nil {
		return err
	}

	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if sc.Err() != nil {
			return err
		}

		entry := sc.Text()
		key, value, exist := strings.Cut(entry, ":")
		if !exist {
			continue
		}

		value = cleanValue(value)

		if len(value) == 0 {
			continue
		}

		switch strings.TrimSpace(key) {
		case "Uid":
			p.Owner = getOwner(value)
		case "Gid":
			p.Group = getGroup(value)
		case "Name":
			p.Name = value
		case "VmSize":
			p.VirtualMem = parseMemoryValue(value)
		case "VmRSS":
			p.ResidentMem = parseMemoryValue(value)
		case "RssFile", "RssShmem":
			p.SharedMem += parseMemoryValue(value)
		}
	}

	return nil
}

func (p *Process) parseStat() error {
	stat, err := os.ReadFile(path.Join("/proc", p.ID, "stat"))
	if err != nil {
		return err
	}

	data := bytes.Fields(stat)
	for i, v := range data {
		switch i {
		case 2:
			p.Status = v[0]
		case 17:
			if pr, err := strconv.Atoi(string(v)); err == nil {
				p.Priority = pr
			}
		case 18:
			if nc, err := strconv.Atoi(string(v)); err == nil {
				p.Nice = nc
			}
		}
	}

	return nil
}

func isPID(dir string, maxPID int) bool {
	base := path.Base(dir)

	pid, err := strconv.Atoi(base)
	if err != nil {
		return false
	}

	if pid > maxPID {
		return false
	}

	return true
}

func getCommand(pid string) string {
	cmd, err := os.ReadFile(path.Join("/proc", pid, "cmdline"))
	if err != nil {
		return ""
	}

	return cleanValue(string(cmd))
}

func getOwner(statusUid string) string {
	uidItems := strings.Fields(statusUid)
	user, err := user.LookupId(uidItems[1])
	if err != nil {
		return ""
	}

	return user.Username
}

func getGroup(statusGuid string) string {
	gidItems := strings.Fields(statusGuid)
	group, err := user.LookupGroupId(gidItems[1])
	if err != nil {
		return ""
	}

	return group.Name
}

func getMaxPID() int {
	p := "/proc/sys/kernel/pid_max"

	max, err := os.ReadFile(p)
	if err != nil {
		return 0
	}

	max = bytes.TrimSpace(max)

	maxPID, err := strconv.Atoi(string(max))
	if err != nil {
		return 0
	}

	return maxPID
}

func getProcesses() ([]Process, error) {
	p := "/proc/"
	dirs, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var processes []Process

	for _, d := range dirs {
		pid := d.Name()
		if isPID(pid, getMaxPID()) {
			p := Process{
				ID:      pid,
				Command: getCommand(pid),
			}

			p.parseStatus()
			p.parseStat()
			processes = append(processes, p)
		}
	}

	return processes, nil
}
