package cmd

import (
	"fmt"
	"github.com/galqiwi/garc/internal/utils/shell"
	"github.com/galqiwi/garc/internal/verify/verifiers/verifier"
	"os/exec"
)

var NotifySendBinary = "notify-send"
var NoNotifyBinaryErr = fmt.Errorf(
	"%v binary is not found, can't send notification",
	NotifySendBinary,
)

func doDisplayNotification(messages []verifier.ErrorMessage) error {
	if ok, err := shell.CommandExists(NotifySendBinary); err != nil || !ok {
		return NoNotifyBinaryErr
	}

	nErrors := getNErrors(messages)
	nWarnings := getNWarns(messages)

	level := "critical"
	notifText := fmt.Sprintf(
		"You have %v errors and %v warnings.",
		nErrors,
		nWarnings,
	)

	if nErrors == 0 {
		level = "normal"
		notifText = fmt.Sprintf(
			"You have %v warnings.",
			nWarnings,
		)
	}

	notifText = fmt.Sprintf("garc verify:\n\n%v", notifText)

	args := []string{notifText}
	args = append(args, "-u", level)

	err := exec.Command(NotifySendBinary, args...).Run()
	if err != nil {
		return err
	}

	return nil
}
