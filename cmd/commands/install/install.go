package install

import (
	"github.com/v3-nu/vv/cmd/utils"
	"github.com/spf13/cobra"
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
