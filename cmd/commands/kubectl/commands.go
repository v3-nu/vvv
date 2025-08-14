package kubectl

import (
	"bytes"
	"fmt"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/v3-nu/vvv/cmd/utils"
	"github.com/v3-nu/vvv/internal/types"
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
			Run:     utils.AliasCommand("kubectl config use-context %s"),
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
			Run:     utils.AliasCommand("kubectl config set-context --current --namespace %s"),
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

func KubectlGetDecodedSecret(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "getsecret [namespace] [secret]",
			Aliases: []string{"getsec", "gsec"},
			Short:   "Get a secret and decode it",
			Args:    cobra.MaximumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) == 2 {
					utils.RunBash("kubectl get secret -n %s %s -o go-template='{{range $k,$v := .data}}{{printf \"%%s: \" $k}}{{if not $v}}{{$v}}{{else}}{{$v | base64decode}}{{end}}{{\"\\n\"}}{{end}}'", args[0], args[1])
				}
				if len(args) == 1 {
					utils.RunBash("kubectl get secret %s -o go-template='{{range $k,$v := .data}}{{printf \"%%s: \" $k}}{{if not $v}}{{$v}}{{else}}{{$v | base64decode}}{{end}}{{\"\\n\"}}{{end}}'", args[0])
				}
			},
		},
	)
}

func KubectlChooseNamespace(cmd *cobra.Command) {
	cmdx := &cobra.Command{
		Use:     "namespace",
		Aliases: []string{"ns", "n"},
		Short:   "Choose a kubeconfig namespace",
		Run: func(cmd *cobra.Command, args []string) {
			app := tview.NewApplication()

			list := tview.NewList()

			namespaces, err := utils.RunBashReturn("kubectl get namespaces -o name")
			if err != nil {
				fmt.Println("Error getting namespaces:", err)
				return
			}

			for _, ns := range bytes.Split(namespaces, []byte{'\n'}) {
				nsx := bytes.TrimPrefix(ns, []byte("namespace/"))
				if len(nsx) == 0 {
					continue
				}
				list = list.AddItem(string(nsx), "", 0, func() {
					utils.RunBash("kubectl config set-context --current --namespace " + string(nsx))
					app.Stop()
				})
			}

			list = list.AddItem("Exit", "Press to exit", 'q', func() {
				app.Stop()
			})

			if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
				panic(err)
			}
		},
	}

	cmd.AddCommand(cmdx)
}

func KubectlChooseContext(cmd *cobra.Command) {
	cmdx := &cobra.Command{
		Use:     "context",
		Aliases: []string{"ctx", "c"},
		Short:   "Choose a kubeconfig context",
		Run: func(cmd *cobra.Command, args []string) {
			app := tview.NewApplication()

			list := tview.NewList()

			kubeconfig, err := types.TryGetKubeConfig()
			if err != nil {
				fmt.Println("Error reading kubeconfig:", err)
				return
			}

			for _, cx := range kubeconfig.Contexts {
				ctxName := cx.Name
				ctxDesc := cx.Context.Namespace

				if kubeconfig.ClustersMap[cx.Context.Cluster].Name != "" {
					ctxDesc = fmt.Sprintf("%s (ns: %s)", kubeconfig.ClustersMap[cx.Context.Cluster].Cluster.Server, cx.Context.Namespace)
				}

				list = list.AddItem(ctxName, ctxDesc, 0, func() {
					utils.RunBash("kubectl config use-context " + cx.Name)
					app.Stop()
				})
			}

			list = list.AddItem("Exit", "Press to exit", 'q', func() {
				app.Stop()
			})

			if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
				panic(err)
			}
		},
	}

	cmd.AddCommand(cmdx)
}
