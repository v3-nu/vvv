package config

import (
	"github.com/clysec/clycli/cmd/utils"
	"github.com/spf13/cobra"
)

var ExportCommands = utils.CommandGroup{
	Command: &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Manage configuration, aliases and settings for clycli",
	},
	Children: []func(*cobra.Command){
		SetConfig,
		GetConfig,
	},
}
