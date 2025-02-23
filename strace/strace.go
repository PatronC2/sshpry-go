package strace

import (
	"bytes"
	"os/exec"
)

type STrace struct {
	cmd *exec.Cmd

	Stdout bytes.Buffer
	Stderr bytes.Buffer

	Flags map[string]string
}

func (s *STrace) Trace() error {
	flags := []string{}

	for k, v := range s.Flags {
		flags = append(flags, k, v)
	}

	s.cmd = exec.Command("strace", flags...)

	s.cmd.Stdout = &s.Stdout
	s.cmd.Stderr = &s.Stderr

	return s.cmd.Start()
}
