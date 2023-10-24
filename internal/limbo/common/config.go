package common

import (
	"github.com/galqiwi/garc/internal/utils/ssh"
)

type LimboConfig struct {
	Client   ssh.ClientConfig
	Hostname string
	Username string
	Port     int
	Path     string
}
