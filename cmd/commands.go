package cmd

import (
	"os"

	"github.com/clysec/clytool/cmd/commands"
	gofigure "github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var rootCommand *cobra.Command

func Execute() {
	figure := gofigure.NewColorFigure("clytool", "doh", "green", true)

	rootCommand = &cobra.Command{
		Use:     "cly",
		Short:   "A command-line utility with various functionality",
		Long:    figure.ColorString(),
		Aliases: []string{"cy"},
	}

	AddCommands(rootCommand)

	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func AddCommands(rootCommand *cobra.Command) {
	commands.AddKubectlListContexts(rootCommand)
	commands.AddKubectlSetContext(rootCommand)
	commands.AddKubectlSetNamespace(rootCommand)
	commands.AddKubectlRemoveFinalizers(rootCommand)
}
