package install

import (
	"github.com/spf13/cobra"
	"github.com/v3-nu/vvv/cmd/utils"
)

var ExportCommands = utils.CommandGroup{
	Command: &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Install software and other things",
	},
	Children: []func(*cobra.Command){
		InstallAzureCli,
		InstallGo,
		InstallHelm,
		InstallKubectl,
		InstallJupyterBashKernel,
		InstallNodejs,
		InstallAcmeShell,
	},
}
