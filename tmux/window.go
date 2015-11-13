package tmux

func NewWindow(session, cmd string) error {
	var err error
	if HasSession(session) {
		_, err = raw("new-window", "-d", "-t", session, cmd)
	} else {
		err = NewSession(session, cmd)
	}
	return err
}
