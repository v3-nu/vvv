package crypto

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/clysec/clycli/cmd/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func PfxToPem(cmd *cobra.Command) {
	command := &cobra.Command{
		Use:     "pfx-to-pem [path-to-pfx-file]",
		Aliases: []string{"pfx2pem", "p2p"},
		Short:   "Convert a pfx file to a set of PEM files (.key, .nopass.key, .crt, .ca.crt, .fullchain.pem)",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]

			fmt.Println("Enter password:")
			password, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println("Error reading password:", err)
				os.Exit(1)
			}

			filename := strings.TrimSuffix(strings.TrimSuffix(filepath.Base(path), ".pfx"), ".p12")
			newFilepath := filepath.Join(filepath.Dir(path), filename)

			nodes := ""
			nodesFlag, err := cmd.Flags().GetBool("nodes")
			if err == nil && nodesFlag {
				confirmFlag, err := cmd.Flags().GetBool("confirm")
				if err != nil || !confirmFlag {
					fmt.Println("Confirm you want to save the private keys without a password")
					res, err := bufio.NewReader(os.Stdin).ReadString('\n')
					if err != nil {
						fmt.Println("Error reading confirmation:", err)
						os.Exit(1)
					}
					if !strings.HasPrefix(strings.ToLower(res), "y") {
						fmt.Println("Exiting...")
						os.Exit(0)
					}
				}

				nodes = "-nodes"
			}

			utils.RunBash(`openssl pkcs12 -in "%s" -out "%s" -nocerts %s -passin pass:%s -passout pass:%s`, path, newFilepath+".key", nodes, password, password)
			utils.RunBash(`openssl pkcs12 -in "%s" -out "%s" %s -passin pass:%s -passout pass:%s`, path, newFilepath+".combined.pem", nodes, password, password)
			utils.RunBash(`openssl pkcs12 -in "%s" -out "%s" -nokeys -clcerts -passin pass:%s`, path, newFilepath+".crt", password)
			utils.RunBash(`openssl pkcs12 -in "%s" -out "%s" -nokeys -cacerts -passin pass:%s`, path, newFilepath+".ca.crt", password)
			utils.RunBash(`openssl pkcs12 -in "%s" -out "%s" -nokeys -passin pass:%s`, path, newFilepath+".fullchain.pem", password)
		},
	}

	command.Flags().BoolP("nodes", "n", false, "Include this flag if you want the resulting private keys to be saved without a password")
	command.Flags().BoolP("confirm", "c", false, "Include this flag to not prompt about saving the private keys without a password")

	cmd.AddCommand(
		command,
	)
}
