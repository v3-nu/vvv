package commands

import (
	"github.com/clysec/clytool/cmd/utils"
	"github.com/spf13/cobra"
)

func AddKubectlListContexts(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "kcontexts",
			Aliases: []string{"kctx", "kcon", "getc"},
			Short:   "List all kubeconfig contexts",
			Run:     utils.AliasCommand("kubectl config get-contexts"),
		},
	)
}

func AddKubectlSetContext(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "ksetcontext [context]",
			Aliases: []string{"ksetc", "ksetctx", "kset", "kcon"},
			Short:   "Set the current-context in the kubeconfig",
			Args:    cobra.ExactArgs(1),
			Run:     utils.AliasCommand("kubectl config set-context"),
		},
	)
}

func AddKubectlSetNamespace(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "ksetnamespace [namespace]",
			Aliases: []string{"ksetns", "ksetn", "kns"},
			Short:   "Set the current namespace in the kubeconfig",
			Args:    cobra.ExactArgs(1),
			Run:     utils.AliasCommand("kubectl config set-context --current --namespace"),
		},
	)
}

func AddKubectlRemoveFinalizers(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:       "krmfinalizers [namespace] [resource]",
			Aliases:   []string{"krmfin", "kdelfinal"},
			Short:     "Remove finalizers from a resource",
			Args:      cobra.ExactArgs(2),
			ValidArgs: []string{"namespace", "resource"},
			Run:       utils.SprintfCommand("kubectl -n %s get %s -o json | jq 'del(.finalizers[])' | kubectl -n %s replace --raw %s -f -"),
		},
	)
}
