package ssh_utils

import (
	"errors"
	"regexp"
)

var InvalidPathErr = errors.New("invalid path")
var pathRegex = regexp.MustCompile(`^/[.\-a-zA-Z0-9/]*$`)

func validatePath(path string) error {
	matched := pathRegex.MatchString(path)
	if !matched {
		return InvalidPathErr
	}
	return nil
}
