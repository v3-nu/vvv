package cmd

import (
	"context"
	"log"
	"os"

	"github.com/v3-nu/vv/cmd/commands/alias"
	"github.com/v3-nu/vv/cmd/commands/crypto"
	"github.com/v3-nu/vv/cmd/commands/install"
	"github.com/v3-nu/vv/cmd/commands/kubectl"
	"github.com/v3-nu/vv/cmd/commands/packages"
	"github.com/v3-nu/vv/cmd/commands/txt"
	"github.com/v3-nu/vv/cmd/commands/uploads"
	"github.com/v3-nu/vv/cmd/utils"
	"github.com/v3-nu/vv/config"
	gofigure "github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var rootCommand *cobra.Command

var registerGroups = []utils.CommandGroup{
	kubectl.ExportCommands,
	install.ExportCommands,
	uploads.ExportCommands,
	packages.ExportCommands,
	config.ExportCommands,
	alias.ExportCommands,
	crypto.ExportCommands,
	txt.ExportCommands,
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
