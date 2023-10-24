package common

import (
	"os"
	"path"
	"strings"
)

var LimboIgnoreFileName = ".limboignore"

func GetLimboIgnore(dir string) []string {
	limboIgnoreData, err := os.ReadFile(path.Join(dir, LimboIgnoreFileName))
	if err != nil {
		return nil
	}
	var output []string
	for _, line := range strings.Split(string(limboIgnoreData), "\n") {
		output = append(output, line)
	}
	return output
}
