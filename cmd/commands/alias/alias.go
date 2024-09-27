package alias

import (
	"fmt"
	"strings"

	"log"

	"github.com/clysec/clycli/cmd/utils"
	"github.com/clysec/clycli/config"
	"github.com/spf13/cobra"
)

var ExportCommands = utils.CommandGroup{
	Command: &cobra.Command{
		Use:   "alias",
		Short: "Manage shell aliases",
		Long:  "Manage shell aliases through in the .clycli/.aliasrc file",
	},
	Children: []func(*cobra.Command){
		ListAliases,
		SetAlias,
	},
}

func InstallAlias(cmd *cobra.Command) {
	// AliasRcRegex := regexp.MustCompile(`(?m)^`)

	cmd.AddCommand(
		&cobra.Command{
			Use:     "install-alias",
			Aliases: []string{"install"},
			Short:   "Install aliases into your shell",
			Run: func(cmd *cobra.Command, args []string) {
				// homedir, err := os.UserHomeDir()
				// if err != nil {
				// 	log.Fatalf("Failed to get user home directory: %v", err)
				// }

				// aliasRcLocation := config.GetAliasRcLocation()

				// bashrcContent, err := os.ReadFile(homedir + "/.bashrc")
				// if err != nil {
				// 	log.Fatalf("Failed to read .bashrc: %v", err)
				// }

			},
		},
	)
}

func ListAliases(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
			Short:   "List aliases",
			Run: func(cmd *cobra.Command, args []string) {
				cfg := cmd.Context().Value(config.ConfigKey("config")).(*config.Config)

				for alias, command := range cfg.Aliases {
					fmt.Printf("%s: %s\n", alias, command)
				}
			},
		},
	)
}

func SetAlias(cmd *cobra.Command) {
	command := &cobra.Command{
		Use:   "set",
		Short: "Set an alias",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := cmd.Context().Value(config.ConfigKey("config")).(*config.Config)

			overwrite, _ := cmd.Flags().GetBool("overwrite")

			if overwrite || cfg.Aliases[args[0]] == "" {
				cfg.Aliases[args[0]] = strings.Join(args[1:], " ")
			} else {
				log.Fatalf("Alias %s already exists, and overwrite is not specified", args[0])
			}

			err := cfg.SaveConfig()
			if err != nil {
				log.Fatalf("Failed to save config: %v", err)
			}
		},
	}

	command.Flags().BoolP("overwrite", "o", false, "Overwrite an existing alias")
	cmd.AddCommand(command)
}
