package config

import (
	"fmt"
	"log"

	"github.com/clysec/clycli/config"
	"github.com/spf13/cobra"
)

func SetConfig(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "set [group] [key] [value]",
			Aliases: []string{"s"},
			Short:   "Set configuration key/value",
			Args:    cobra.ExactArgs(3),
			Run: func(cmd *cobra.Command, args []string) {
				cfg := cmd.Context().Value(config.ConfigKey("config")).(*config.Config)

				err := cfg.SetConfigValue(args[0], args[1], args[2])
				if err != nil {
					log.Fatalf("Failed to set configuration value: %v", err)
				}
			},
		},
	)
}

func GetConfig(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "get [group] [key]",
			Aliases: []string{"g"},
			Short:   "Get configuration value",
			Args:    cobra.MaximumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				cfg := cmd.Context().Value(config.ConfigKey("config")).(*config.Config)
				avail := cfg.GetCurrentSettingsMap()

				fmt.Println("Settings:", "\r\n")

				if len(args) == 2 && avail[args[0]] != nil {
					fmt.Printf("%s.%s = %s\n", args[0], args[1], avail[args[0]][args[1]])
					return
				}

				if len(args) == 1 {
					if avail[args[0]] != nil {
						fmt.Println(args[0] + ":")
						for key := range avail[args[0]] {
							fmt.Printf("  %s = %s\n", key, avail[args[0]][key])
						}

						return
					}
				}

				for group, keys := range avail {
					fmt.Println(group + ":")
					for key, value := range keys {
						fmt.Printf("  %s = %s\n", key, value)
					}
				}
			},
		},
	)
}
