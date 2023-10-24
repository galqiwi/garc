package ssh

type Client interface {
	Host(hostname string, port int) Host
}

type ClientConfig struct {
	KeyPath string
}

type client struct {
	cfg *ClientConfig
}

func (c *client) Host(hostname string, port int) Host {
	return &host{
		client:   c,
		hostname: hostname,
		port:     port,
	}
}

func NewClient(cfg *ClientConfig) Client {
	return &client{cfg: cfg}
}

var _ Client = &client{}
