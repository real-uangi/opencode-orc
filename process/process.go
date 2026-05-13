package process

import (
	"fmt"
	"io"
	"os/exec"
)

// Process manages the opencode subprocess
type Process struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

// New creates a new process manager
func New(args []string) (*Process, error) {
	cmd := exec.Command("opencode", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	return &Process{
		cmd:    cmd,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

// Start starts the subprocess
func (p *Process) Start() error {
	return p.cmd.Start()
}

// Wait waits for the subprocess to finish
func (p *Process) Wait() error {
	return p.cmd.Wait()
}

// Stdout returns the stdout reader
func (p *Process) Stdout() io.ReadCloser {
	return p.stdout
}

// Stderr returns the stderr reader
func (p *Process) Stderr() io.ReadCloser {
	return p.stderr
}

// ExitCode returns the exit code of the process
func (p *Process) ExitCode() int {
	if p.cmd.ProcessState != nil {
		return p.cmd.ProcessState.ExitCode()
	}
	return -1
}
