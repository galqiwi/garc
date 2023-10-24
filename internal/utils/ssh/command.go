package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

type Command interface {
	Run() error
	SetStderr(stderr io.Writer)
	SetStdout(stdout io.Writer)
	SetStdin(stdin io.Reader)
}

type command struct {
	host     *host
	username string
	cmd      string
	stderr   io.Writer
	stdout   io.Writer
	stdin    io.Reader
}

func (c *command) SetStderr(stderr io.Writer) {
	c.stderr = stderr
}

func (c *command) SetStdout(stdout io.Writer) {
	c.stdout = stdout
}

func (c *command) SetStdin(stdin io.Reader) {
	c.stdin = stdin
}

func (c *command) Run() error {
	key, err := os.ReadFile(c.host.client.cfg.KeyPath)
	if err != nil {
		return fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: c.username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial(
		"tcp",
		fmt.Sprintf("%v:%v", c.host.hostname, c.host.port),
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	session.Stdout = c.stdout
	session.Stderr = c.stderr
	session.Stdin = c.stdin

	return session.Run(c.cmd)
}

var _ Command = &command{}
