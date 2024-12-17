package kubectl

import (
	"github.com/spf13/cobra"
	"github.com/v3-nu/vvv/cmd/utils"
)

var ExportCommands = utils.CommandGroup{
	Command: &cobra.Command{
		Use:     "kubectl [command] [flags]",
		Short:   "kubectl",
		Long:    "kubectl",
		Aliases: []string{"k"},
	},
	Children: []func(*cobra.Command){
		KubectlListContexts,
		KubectlSetContext,
		KubectlSetNamespace,
		KubectlRemoveFinalizers,
		KubectlGetDecodedSecret,
	},
}
