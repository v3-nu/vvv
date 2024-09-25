package cmd

import (
	"context"
	"log"
	"os"

	configcmd "github.com/clysec/clycli/cmd/commands/config"
	"github.com/clysec/clycli/cmd/commands/install"
	"github.com/clysec/clycli/cmd/commands/kubectl"
	"github.com/clysec/clycli/cmd/commands/packages"
	"github.com/clysec/clycli/cmd/commands/uploads"
	"github.com/clysec/clycli/cmd/utils"
	"github.com/clysec/clycli/config"
	gofigure "github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var rootCommand *cobra.Command

var registerGroups = []utils.CommandGroup{
	kubectl.ExportCommands,
	install.ExportCommands,
	uploads.ExportCommands,
	packages.ExportCommands,
	configcmd.ExportCommands,
}

func Execute() {
	figure := gofigure.NewColorFigure("clycli", "doh", "green", true)

	rootCommand = &cobra.Command{
		Use:     "clycli",
		Short:   "A command-line utility with various functionality",
		Long:    figure.ColorString(),
		Aliases: []string{"cly", "cy", "cc"},
		// PreRun: func(cmd *cobra.Command, args []string) {
		// 	log.Fatalf("ASDASD")
		// },
	}

	for _, group := range registerGroups {
		group.Register(rootCommand)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
	}

	ctx := context.WithValue(context.Background(), config.ConfigKey("config"), cfg)
	err = rootCommand.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
