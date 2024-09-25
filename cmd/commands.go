package cmd

import (
	"os"

	"github.com/clysec/clycli/cmd/commands/install"
	"github.com/clysec/clycli/cmd/commands/kubectl"
	"github.com/clysec/clycli/cmd/commands/packages"
	"github.com/clysec/clycli/cmd/commands/uploads"
	"github.com/clysec/clycli/cmd/utils"
	gofigure "github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var rootCommand *cobra.Command

var registerGroups = []utils.CommandGroup{
	kubectl.ExportCommands,
	install.ExportCommands,
	uploads.ExportCommands,
	packages.ExportCommands,
}

func Execute() {
	figure := gofigure.NewColorFigure("clycli", "doh", "green", true)

	rootCommand = &cobra.Command{
		Use:     "clycli",
		Short:   "A command-line utility with various functionality",
		Long:    figure.ColorString(),
		Aliases: []string{"cly", "cy", "cc"},
	}

	for _, group := range registerGroups {
		group.Register(rootCommand)
	}

	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
