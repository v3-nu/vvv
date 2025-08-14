package secenv

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/v3-nu/vvv/cmd/utils"
)

var ExportCommands utils.CommandGroup = func() utils.CommandGroup {
	command := utils.CommandGroup{
		Command: &cobra.Command{
			Use:     "secenv",
			Aliases: []string{"env", "se"},
			Short:   "Secure environment variables, files and related commands",
		},
		Children: []func(*cobra.Command){
			Init,
			ListEnv,
		},
	}

	command.Command.PersistentFlags().StringP("env", "e", "vvv", "Environment to use. Default is vvv")

	return command
}()

func Init(cmd *cobra.Command) {
	comm := &cobra.Command{
		Use:     "init",
		Short:   "Initialize the secure environment",
		Aliases: []string{"golang"},
		Run: func(cmd *cobra.Command, args []string) {
			env := cmd.Flag("env").Value.String()
			if env == "" {
				env = "vvv"
			}

			desc := cmd.Flag("description").Value.String()
			if desc == "" {
				desc = utils.AskString("Enter a description for the secure environment", false)
			}

			secEnv := &SecureEnvironment{
				Name:  env,
				Desc:  desc,
				Vars:  make(map[string]SecureItem),
				Files: make(map[string]SecureItem),
			}

			secEnv.Init()

			fmt.Println("Secure environment initialized successfully.")
			fmt.Printf("Name: %s\n", secEnv.Name)
			fmt.Printf("Description: %s\n", secEnv.Desc)
		},
	}

	comm.Flags().StringP("description", "d", "", "Set a description for the secure environment")

	cmd.AddCommand(comm)
}

func ListEnv(cmd *cobra.Command) {
	cmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all secure environments",
		Run: func(cmd *cobra.Command, args []string) {
			env := cmd.Flag("env").Value.String()
			if env == "" {
				env = "vvv"
			}

			secEnv := &SecureEnvironment{}
			lst, err := secEnv.List()
			if err != nil {
				fmt.Println("Error listing secure environments:", err)
				return
			}

			fmt.Println("Environments: ")

			for _, name := range lst {
				fmt.Println("[" + name + "]")
			}
		},
	})
}
