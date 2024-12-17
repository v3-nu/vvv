package uploads

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/v3-nu/vvv/cmd/utils"
)

var ExportCommands = utils.CommandGroup{
	Command: &cobra.Command{
		Use:     "uploads",
		Aliases: []string{"upload", "u"},
		Short:   "Upload files and/or data to various services",
	},
	Children: []func(*cobra.Command){
		UploadPastebin,
	},
}

func UploadPastebin(cmd *cobra.Command) {
	cmd.AddCommand(
		&cobra.Command{
			Use:     "pastebin",
			Aliases: []string{"paste", "pb"},
			Short:   "Upload a file to pastebin",
			Long:    "Upload a file to pastebin. To configure the pastebin server and type of service, see the \"vvv config upload pastebin\" command.",
			Run: func(cmd *cobra.Command, args []string) {
				conn, err := net.Dial("tcp", "termbin.com:9999")
				if err != nil {
					fmt.Println("Error connecting to termbin: ", err)
				}

				io.Copy(conn, os.Stdin)

				data, err := io.ReadAll(conn)
				if err != nil {
					fmt.Println("Error reading from termbin: ", err)
				}

				fmt.Println(string(data))

				conn.Close()
			},
		},
	)
}
