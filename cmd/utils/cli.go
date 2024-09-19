package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func StringSliceToAnySlice(args []string) []any {
	anyargs := make([]any, len(args))

	for i, v := range args {
		anyargs[i] = v
	}

	return anyargs
}

func AliasCommand(command string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		RunBash(command, StringSliceToAnySlice(args)...)
	}
}

func AliasCommandArgpos(command string, argpos ...int) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		anyargs := make([]any, len(argpos))

		for i, v := range argpos {
			anyargs[i] = args[v]
		}

		RunBash(command, anyargs...)
	}
}

// func AliasCommand(command string) func(cmd *cobra.Command, args []string) {
// 	return func(cmd *cobra.Command, args []string) {
// 		SplitAndRun(command, " ", args...)
// 	}
// }

// func AliasCommandBash(command string) func(cmd *cobra.Command, args []string) {
// 	return func(cmd *cobra.Command, args []string) {
// 		RunBash(command)
// 	}
// }

// func SprintfCommand(command string, args ...string) func(cmd *cobra.Command, args []string) {
// 	return func(cmd *cobra.Command, args []string) {
// 		anyargs := make([]any, len(args))

// 		for i, v := range args {
// 			anyargs[i] = v
// 		}

// 		SplitAndRun(fmt.Sprintf(command, anyargs...), " ")
// 	}
// }

// func SprintfCommandBash(command string, args ...an) func(cmd *cobra.Command, args []string) {
// 	return func(cmd *cobra.Command, args []string) {
// 		RunBash(fmt.Sprintf(command, args...))
// 	}
// }

func RunBash(command string, args ...any) {
	cx := exec.Command("bash", "-c", fmt.Sprintf(command, args...))
	cx.Env = os.Environ()
	cx.Stdout = os.Stdout
	cx.Stderr = os.Stderr

	err := cx.Run()
	if err != nil {
		fmt.Println("Error running bash:", err)
		os.Exit(cx.ProcessState.ExitCode())
	}
}

// func SplitAndRun(command, separator string, args ...string) {
// 	split := strings.Split(command, " ")

// 	cx := exec.Command(split[0], append(split[1:], args...)...)
// 	cx.Env = os.Environ()

// 	cx.Stdout = os.Stdout
// 	cx.Stderr = os.Stderr

// 	err := cx.Run()
// 	if err != nil {
// 		fmt.Println("Error listing contexts:", err)
// 		os.Exit(cx.ProcessState.ExitCode())
// 	}
// }
