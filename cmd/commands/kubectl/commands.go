package kubectl

import (
	"github.com/clysec/clycli/cmd/utils"
	"github.com/spf13/cobra"
)

func KubectlListContexts(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "getcontexts",
			Aliases: []string{"getctx", "getc"},
			Short:   "List kubeconfig contexts",
			Run:     utils.AliasCommand("kubectl config get-contexts"),
		},
	)
}

func KubectlSetContext(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "usecontext [context]",
			Aliases: []string{"setc", "con", "usectx", "usec"},
			Short:   "Set the current-context in the kubeconfig",
			Args:    cobra.ExactArgs(1),
			Run:     utils.AliasCommand("kubectl config set-context"),
		},
	)
}

func KubectlSetNamespace(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "setns [namespace]",
			Aliases: []string{"setn", "ns"},
			Short:   "Set the current namespace in the kubeconfig",
			Args:    cobra.ExactArgs(1),
			Run:     utils.AliasCommand("kubectl config set-context --current --namespace"),
		},
	)
}

func KubectlRemoveFinalizers(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:       "patchfinalizezrs [namespace] [resource]",
			Aliases:   []string{"patchfin", "delfin"},
			Short:     "Remove finalizers from a resource",
			Args:      cobra.ExactArgs(2),
			ValidArgs: []string{"namespace", "resource"},
			Run:       utils.AliasCommandArgpos("kubectl -n %s get %s -o json | jq 'del(.finalizers[])' | kubectl -n %s replace --raw %s -f -", 0, 1, 0, 1),
		},
	)
}
