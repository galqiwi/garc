package ssh_utils

import (
	"github.com/galqiwi/garc/internal/utils/ssh"
	"strings"
)

func RunAndGetStdout(command ssh.Command) (string, error) {
	buf := strings.Builder{}
	command.SetStdout(&buf)
	err := command.Run()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
