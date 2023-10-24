package shell

import "os/exec"

func CommandExists(cmd string) (bool, error) {
	err := exec.Command("which", cmd).Run()
	if err == nil {
		return true, nil
	}
	if _, exited := err.(*exec.ExitError); exited {
		return false, nil
	}
	return false, err
}
