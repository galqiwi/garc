package ssh

type Host interface {
	Command(username string, cmd string) Command
}

type host struct {
	client   *client
	hostname string
	port     int
}

func (h *host) Command(username string, cmd string) Command {
	return &command{
		host:     h,
		username: username,
		cmd:      cmd,
	}
}

var _ Host = &host{}
