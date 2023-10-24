package misc

import "os"

func GetCurrentExecutablePath() (string, error) {
	return os.Readlink("/proc/self/exe")
}
