package tmux

import (
	"fmt"
	"strings"
)

func NewWindow(session, name, cmd string) error {
	var err error
	if HasSession(session) {
		_, err = raw("new-window", "-d", "-n", name, "-t", session, cmd)
	} else {
		err = NewSession(session, name, cmd)
	}
	return err
}

func KillWindow(session, name string) error {
	_, err := raw("kill-window", "-t", fmt.Sprintf("%s:%s", session, name))
	return err
}

func ListWindows(session string) ([]string, error) {
	out, err := raw("list-windows", "-t", session)

	split := strings.Split(out, "\n")

	windows := []string{}
	for _, w := range split {
		if w != "" {
			windows = append(windows, w)
		}
	}
	return windows, err
}
