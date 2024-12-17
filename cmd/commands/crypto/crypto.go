package crypto

import (
	"github.com/spf13/cobra"
	"github.com/v3-nu/vvv/cmd/utils"
)

var ExportCommands = utils.CommandGroup{
	Command: &cobra.Command{
		Use:     "crypto",
		Aliases: []string{"ssl", "crypt", "tls"},
		Short:   "Crypto related commands",
	},
	Children: []func(*cobra.Command){
		PfxToPem,
	},
}
