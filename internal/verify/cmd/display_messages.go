package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/galqiwi/garc/internal/verify/verifiers/verifier"
	"github.com/rodaine/table"
	"strings"
)

func getNErrors(messages []verifier.ErrorMessage) int {
	nErrors := 0

	for _, message := range messages {
		if message.Level == verifier.Error {
			nErrors += 1
		}
	}

	return nErrors
}

func getNWarns(messages []verifier.ErrorMessage) int {
	nWarns := 0

	for _, message := range messages {
		if message.Level == verifier.Warn {
			nWarns += 1
		}
	}

	return nWarns
}

func displayMessages(messages []verifier.ErrorMessage) {
	fmt.Printf(
		"Got %v errors and %v warnings.\n",
		getNErrors(messages),
		getNWarns(messages),
	)

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()

	tbl := table.New("Level", "Path", "Reason")
	tbl.WithHeaderFormatter(headerFmt)
	tbl.WithFirstColumnFormatter(func(s string, i ...interface{}) string {
		SColor := color.FgYellow
		if strings.Contains(i[0].(string), "Error") {
			SColor = color.FgRed
		}

		return color.New(SColor).SprintfFunc()(s, i...)
	})

	for _, message := range messages {
		tbl.AddRow(verifier.MessageLevelMap[message.Level], message.Path, message.Reason)
	}

	tbl.Print()
}
