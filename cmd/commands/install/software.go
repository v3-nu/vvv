package install

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/v3-nu/vv/cmd/utils"
	"github.com/spf13/cobra"
)

func InstallAzureCli(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "azure-cli",
			Aliases: []string{"azurecli"},
			Short:   "Install Azure CLI (Debian)",
			Run:     utils.AliasCommand(fmt.Sprintf("curl -sL https://aka.ms/InstallAzureCLIDeb | %s", utils.SudoIfNotRoot("bash"))),
		},
	)
}

func InstallGo(cmd *cobra.Command) {
	goVersion := "1.23.0"
	goArch := strings.Join([]string{runtime.GOOS, runtime.GOARCH}, "-")

	fetchUrl := fmt.Sprintf("https://go.dev/dl/go%s.%s.tar.gz", goVersion, goArch)

	cmd.AddCommand(
		&cobra.Command{
			Use:     "go",
			Short:   "Install Go",
			Aliases: []string{"golang"},
			Run: func(cmd *cobra.Command, args []string) {
				utils.RunBash("wget %s -O /tmp/go.tar.gz && %s && %s && rm /tmp/go.tar.gz", fetchUrl, utils.SudoIfNotRoot("rm -rf /usr/local/go"), utils.SudoIfNotRoot("tar -C /usr/local -xzf /tmp/go.tar.gz"))

				if strings.Contains(os.Getenv("PATH"), "/usr/local/go/bin") {
					os.Exit(0)
				}

				if runtime.GOOS == "linux" {
					bashrc, err := os.ReadFile("~/.bashrc")
					if err != nil {
						fmt.Println("Error reading .bashrc:", err)
						os.Exit(1)
					}

					if !strings.Contains(string(bashrc), "/usr/local/go/bin") {
						bashrc, err := os.OpenFile("~/.bashrc", os.O_RDWR, 0o644)
						if err != nil {
							fmt.Println("Error opening .bashrc:", err)
							os.Exit(1)
						}

						bashrc.WriteString("export PATH=$PATH:/usr/local/go/bin")
						bashrc.Close()
					}
				}
			},
		},
	)
}

func InstallHelm(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:   "helm",
			Short: "Install Helm",
			Run:   utils.AliasCommand("wget https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 -O /tmp/install_helm_3.sh && chmod +x /tmp/install_helm_3.sh && /tmp/install_helm_3.sh"),
		},
	)
}

func InstallJupyterBashKernel(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:   "jupyter-bash-kernel",
			Short: "Install Jupyter Bash Kernel",
			Run:   utils.AliasCommand("pip3 install --upgrade bash_kernel && python3 -m bash_kernel.install"),
		},
	)
}

func InstallKubectl(cmd *cobra.Command) {
	commands := []string{
		"curl -fsSL \"https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl\" -o /tmp/v3-kubectl",
		utils.SudoIfNotRoot("install -o root -g root -m 0755 /tmp/v3-kubectl /usr/local/bin/kubectl"),
		utils.SudoIfNotRoot("bash -c 'kubectl completion bash > /etc/bash_completion.d/kubectl'"),
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "kubectl",
			Short: "Install Kubectl",
			Run:   utils.AliasCommand(strings.Join(commands, " && ")),
		},
	)
}

func InstallNodejs(cmd *cobra.Command) {
	nodeVersion := "21.x"

	cmd.AddCommand(
		&cobra.Command{
			Use:   "nodejs",
			Short: "Install NodeJS",
			Run:   utils.AliasCommand(fmt.Sprintf("curl -fsSL https://deb.nodesource.com/setup_%s | sudo -E bash - && sudo apt-get install -y nodejs npm", nodeVersion)),
		},
	)
}

func InstallAcmeShell(cmd *cobra.Command) {
	acmeRepo := "https://github.com/acmesh-official/acme.sh.git"

	acmeCommand := &cobra.Command{
		Use:     "acmesh",
		Aliases: []string{"acme.sh"},
		Short:   "Install Acme Shell (acme.sh)",
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure directories exist
			utils.RunBash(utils.SudoIfNotRoot("mkdir -p %s %s && chown -R %s:%s %s %s"), cmd.Flag("installDirectory").Value.String(), cmd.Flag("configDirectory").Value.String(), cmd.Flag("user").Value.String(), cmd.Flag("user").Value.String(), cmd.Flag("installDirectory").Value.String(), cmd.Flag("configDirectory").Value.String())

			// Install acme.sh
			utils.RunBash(utils.SudoIfNotRoot("git clone %s %s && cd %s && ./acme.sh --install --home %s --config-home %s --user %s"), acmeRepo, cmd.Flag("installDirectory").Value.String(), cmd.Flag("installDirectory").Value.String(), cmd.Flag("installDirectory").Value.String(), cmd.Flag("configDirectory").Value.String(), cmd.Flag("user").Value.String())

			// Add to path
			if !cmd.Flag("skip-alias").Changed {
				utils.RunBash(utils.SudoIfNotRoot("ln -s %s/acme.sh /usr/local/bin/acme.sh"), cmd.Flag("installDirectory").Value.String())
			}

			// If email is provided, register for an account
			if cmd.Flag("email").Changed {
				utils.RunBash(utils.SudoIfNotRoot("acme.sh --register-account -m %s"), cmd.Flag("email").Value.String())
			}

			fmt.Println("Acme.sh installed successfully")
		},
	}

	acmeCommand.Flags().StringP("email", "e", "", "Email address to use for ACME account")
	acmeCommand.Flags().StringP("acme-server", "s", "https://acme.zerossl.com/v2/DV90", "ACME server to use, default is zerossl")

	acmeCommand.Flags().StringP("installDirectory", "i", "/var/lib/acme.sh", "The install directory for acme.sh")
	acmeCommand.Flags().StringP("configDirectory", "c", "/etc/ssl/acme.sh", "The install directory for acme.sh")
	acmeCommand.Flags().StringP("user", "u", "root", "The user to run acme.sh as")

	acmeCommand.Flags().BoolP("skip-alias", "a", false, "Skip adding acme.sh to path")

	cmd.AddCommand(acmeCommand)
}
