package utils

import "fmt"

var (
	sudoIfNotRoot = "if [[ $EUID -ne 0 ]]; then sudo %s; else %s; fi"
)

func SudoIfNotRoot(command string) string {
	return fmt.Sprintf(sudoIfNotRoot, command, command)
}
