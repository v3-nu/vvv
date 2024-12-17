package packages

import (
	"log"
	"strings"

	"github.com/v3-nu/vv/cmd/utils"
	"github.com/spf13/cobra"
)

func GetVendor(cmd *cobra.Command, args []string) Vendor {
	backend, err := cmd.Flags().GetString("backend")
	if err != nil {
		log.Fatalf("Error getting backend: %v", err)
	}

	if backend == "auto" {
		backend = BestGuessPackageManager()
	}

	if VendorMap[backend].Name == "" {
		allowedBackends := make([]string, 0, len(VendorMap))
		for _, vendor := range VendorMap {
			allowedBackends = append(allowedBackends, vendor.Name)
		}
		log.Fatalf("Backend %s not supported, allowed backends: %s or auto", backend, strings.Join(allowedBackends, ", "))
	}

	confirmFlag, err := cmd.Flags().GetBool("yes")
	if err != nil {
		log.Fatalf("Error getting confirm flag: %v", err)
	}

	vendor := VendorMap[backend]
	if !confirmFlag {
		vendor.ConfirmOption = ""
	}

	if len(args) > 0 {
		vendor.Packages = strings.Join(args, " ")
	}

	return vendor
}

func RunCommand(name string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		vendor := GetVendor(cmd, args)

		command := ""

		if name == "list-installed" {
			command = vendor.ListInstalledCommand()
		} else if name == "update" {
			command = vendor.UpdateCommand()
		} else if name == "upgrade" && len(args) == 0 || name == "upgrade-all" {
			command = vendor.UpgradeAllCommand()
		} else if name == "upgrade" {
			command = vendor.UpgradeCommand()
		} else if name == "install" {
			command = vendor.InstallCommand()
		} else if name == "remove" {
			command = vendor.RemoveCommand()
		} else if name == "search" {
			command = vendor.SearchCommand()
		} else if name == "info" {
			command = vendor.InfoCommand()
		} else {
			log.Fatalf("Command %s not supported, valid options are: list-installed, update, upgrade, upgrade-all, install, remove, search, info", name)
		}

		command = vendor.TemplateString(command)
		utils.RunBash(utils.SudoIfNotRoot(command))

	}
}

func AddAllCommands(cmd *cobra.Command) {
	commands := []*cobra.Command{
		{
			Use:     "list-installed",
			Aliases: []string{"list", "lsi", "ls"},
			Short:   "List installed packages",
			Run:     RunCommand("list-installed"),
		},
		{
			Use:     "update",
			Aliases: []string{},
			Short:   "Update packages",
			Run:     RunCommand("update"),
		},
		{
			Use:     "upgrade",
			Aliases: []string{},
			Short:   "Upgrade packages",
			Run:     RunCommand("upgrade"),
		},
		{
			Use:     "upgrade-all",
			Aliases: []string{},
			Short:   "Upgrade all packages",
			Run:     RunCommand("upgrade-all"),
		},
		{
			Use:     "install",
			Aliases: []string{"i", "add"},
			Short:   "Install packages",
			Run:     RunCommand("install"),
		},
		{
			Use:     "remove",
			Aliases: []string{"uninstall", "delete", "rm"},
			Short:   "Remove packages",
			Run:     RunCommand("remove"),
		},
		{
			Use:     "search",
			Aliases: []string{"find"},
			Short:   "Search for packages",
			Run:     RunCommand("search"),
		},
		{
			Use:     "info",
			Aliases: []string{"show", "i"},
			Short:   "Get info about packages",
			Run:     RunCommand("info"),
		},
	}

	cmd.AddCommand(commands...)
}

func GetCobraCommandGroup() utils.CommandGroup {
	command := utils.CommandGroup{
		Command: &cobra.Command{
			Use:     "packages [command] [flags]",
			Short:   "packages",
			Long:    "packages",
			Example: "clycli packages install -y <package>",
			Aliases: []string{"package", "pkg", "p"},
		},
		Children: []func(*cobra.Command){
			AddAllCommands,
		},
	}

	command.Command.PersistentFlags().BoolP("yes", "y", false, "Assume yes to all prompts")
	command.Command.PersistentFlags().StringP("backend", "b", "auto", "Backend to use for the command, or auto")

	return command
}

var ExportCommands = GetCobraCommandGroup()
