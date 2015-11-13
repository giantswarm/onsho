package tmux

import (
	"fmt"
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

func KillWindow(name string) error {
	_, err := raw("kill-window", "-t", name)
	return err
}

func ListWindows(session string) ([]string, error) {
	windows, err := raw("list-windows", "-t", session)
	fmt.Println(windows)
	return []string{}, err
}
