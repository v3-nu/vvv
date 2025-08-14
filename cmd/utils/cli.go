package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
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
	countArgs := strings.Count(command, "%s")
	if countArgs < len(args) {
		repeat := len(args) - countArgs
		command += strings.Repeat(" %s", repeat)
	}

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

func GetStringFlag(cmd *cobra.Command, name string, defaultValue string) string {
	flag := cmd.Flag(name)
	if flag == nil || flag.Value == nil {
		return defaultValue
	}

	value := flag.Value.String()
	if value == "" {
		return defaultValue
	}

	return value
}

func AskString(prompt string, hidden bool) string {
	fmt.Println(prompt)

	if hidden {
		val, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}

		return string(val)
	}

	reader := bufio.NewReader(os.Stdin)
	val, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	return val
}

// // func SplitAndRun(command, separator string, args ...string) {
// // 	split := strings.Split(command, " ")

// // 	cx := exec.Command(split[0], append(split[1:], args...)...)
// // 	cx.Env = os.Environ()

// // 	cx.Stdout = os.Stdout
// // 	cx.Stderr = os.Stderr

// // 	err := cx.Run()
// // 	if err != nil {
// // 		fmt.Println("Error listing contexts:", err)
// // 		os.Exit(cx.ProcessState.ExitCode())
// // 	}
// // }

// func AskPassword(prompt string) string {
// 	fmt.Println(prompt)
// 	var inpt string
// 	_, err := fmt.Scanln(&inpt)
// 	if err != nil {
// 		fmt.Println("Error reading input:", err)
// 		os.Exit(1)
// 	}

// 	return inpt
// }
