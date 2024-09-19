package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func AliasCommand(command string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		SplitAndRun(command, " ", args...)
	}
}

func SprintfCommand(command string, args ...string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		anyargs := make([]any, len(args))

		for i, v := range args {
			anyargs[i] = v
		}

		SplitAndRun(fmt.Sprintf(command, anyargs...), " ")
	}
}

func SplitAndRun(command, separator string, args ...string) {
	split := strings.Split(command, " ")

	cx := exec.Command(split[0], append(split[1:], args...)...)
	cx.Env = os.Environ()

	cx.Stdout = os.Stdout
	cx.Stderr = os.Stderr

	err := cx.Run()
	if err != nil {
		fmt.Println("Error listing contexts:", err)
		os.Exit(cx.ProcessState.ExitCode())
	}
}
