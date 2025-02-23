package main

import (
	"fmt"
	"time"

	"github.com/PatronC2/sshpry-go/strace"
)

func main() {
	s := strace.STrace{
		Flags: map[string]string{
			"-s": "16384",
			"-p": "12",
			"-e": "read,write",
		},
	}

	s.Trace()

	for {
		fmt.Printf("Err: %s\n", s.Stderr.String())
		fmt.Printf("Out: %s\n", s.Stdout.String())

		s.Stderr.Reset()
		s.Stdout.Reset()

		time.Sleep(1 * time.Second)

	}
}
