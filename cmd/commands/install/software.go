package install

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/clysec/clycli/cmd/utils"
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
