package kubectl

import (
	"github.com/clysec/clycli/cmd/utils"
	"github.com/spf13/cobra"
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
	},
}
