package tmux

func KillSession(name string) error {
	if HasSession(name) {
		_, err := raw("kill-session", "-t", name)
		return err
	}
	return nil
}

func HasSession(name string) bool {
	_, err := raw("has-session", "-t", name)
	if err != nil {
		return false
	} else {
		return true
	}
}

func NewSession(session, name, cmd string) error {
	_, err := raw("new-session", "-d", "-n", name, "-s", session, cmd)
	return err
}
