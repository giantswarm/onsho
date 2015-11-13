package tmux

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

var (
	ErrNotFound = errors.New("not found")
)

func raw(subcmd string, args ...string) ([]byte, error) {
	cmd := exec.Command("tmux", append([]string{subcmd}, args...)...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s - %v", stdout.String(), stderr.String(), err)
	}

	return []byte(stdout.String()), nil
}
