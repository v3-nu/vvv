package test

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/v3-nu/vvv/cmd/utils"
)

var ExportCommands utils.CommandGroup = func() utils.CommandGroup {
	command := utils.CommandGroup{
		Command: &cobra.Command{
			Use:     "test",
			Aliases: []string{"tst"},
			Short:   "Test functions for development and debugging",
		},
		Children: []func(*cobra.Command){
			TestTree,
		},
	}

	command.Command.PersistentFlags().StringP("env", "e", "vvv", "Environment to use. Default is vvv")

	return command
}()

func TestTree(cmd *cobra.Command) {
	comm := &cobra.Command{
		Use:     "tree",
		Short:   "Create a test tree of directories and files up to a specified size (default 100MB)",
		Aliases: []string{"tree", "tt"},
		Run: func(cmd *cobra.Command, args []string) {
			dest := cmd.Flag("destination").Value.String()
			if dest == "" {
				dest = "/tmp/vvvtree"
			}

			size, err := cmd.Flags().GetInt("size")
			if err != nil {
				fmt.Printf("Failed to get size flag: %v", err)
				return
			}

			if size <= 0 {
				fmt.Println("Size must be a positive integer")
				return
			}

			err = CreateTestTree(dest, size)
			if err != nil {
				fmt.Printf("Failed to create test tree: %v", err)
				return
			}

			fmt.Println("Successfully created test tree of size", size, "MB at", dest)
		},
	}

	comm.Flags().StringP("destination", "d", "/tmp/vvvtree", "Where to create the test tree")
	comm.Flags().IntP("size", "s", 100, "Maximum size of the test tree in MB (default 100MB)")

	cmd.AddCommand(comm)
}
