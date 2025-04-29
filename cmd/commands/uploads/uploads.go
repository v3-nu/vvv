package uploads

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/clysec/greq"
	"github.com/spf13/cobra"
	"github.com/v3-nu/vv/cmd/utils"
	"github.com/v3-nu/vv/config"
)

var ExportCommands = utils.CommandGroup{
	Command: &cobra.Command{
		Use:     "uploads",
		Aliases: []string{"upload", "u"},
		Short:   "Upload files and/or data to various services",
	},
	Children: []func(*cobra.Command){
		UploadPastebin,
		UploadTransfersh,
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

func UploadTransfersh(cmd *cobra.Command) {
	configs := utils.FindContext(cmd).Value(config.ConfigKey("config")).(*config.Config)

	transfershUrl := configs.Settings.Uploads.TransfershUrl
	transfershUser := configs.Settings.Uploads.TransfershUsername
	transfershPass := configs.Settings.Uploads.TransfershPassword

	transfershCommand := &cobra.Command{
		Use:     "transfersh",
		Aliases: []string{"tsh", "file"},
		Short:   "Upload file to transfersh instance",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			file, err := os.Stat(filename)
			if err != nil {
				filename := args[len(args)-1]
				file, err = os.Stat(filename)
				if err != nil {
					fmt.Println("Error getting file info: ", err)
					return
				}
			}

			if file.IsDir() {
				fmt.Printf("Error: %s is a directory", filename)
				return
			}

			fmt.Println("Uploading file: ", filename)

			hostVar := utils.GetStringFlag(cmd, "url", transfershUrl)
			if hostVar == "" {
				fmt.Println("Error: No transfer.sh host provided")
				return
			}

			userVar := utils.GetStringFlag(cmd, "username", transfershUser)
			passVar := utils.GetStringFlag(cmd, "password", transfershPass)

			stream, err := os.OpenFile(filename, os.O_RDONLY, 0644)
			if err != nil {
				fmt.Println("Error opening file: ", err)
				return
			}

			defer stream.Close()

			if !strings.HasPrefix(hostVar, "http") {
				hostVar = "https://" + hostVar
			}

			filenameVar := cmd.Flag("filename").Value.String()
			if filenameVar == "" {
				filenameVar = file.Name()
			}

			hostVar, err = url.JoinPath(hostVar, filenameVar)
			if err != nil {
				fmt.Println("Error joining URL path: ", err)
				return
			}

			// Upload the file to transfersh
			resp, err := greq.PutRequest(hostVar).WithAuth(&greq.BasicAuth{
				Username: userVar,
				Password: passVar,
			}).
				WithReaderBody(stream).
				Execute()

			if err != nil {
				fmt.Println("Error uploading file: ", err)
				return
			}

			if resp.StatusCode != 200 {
				fmt.Println("Error uploading file: Code ", resp.StatusCode)
				return
			}

			body, err := resp.BodyString()
			if err != nil {
				fmt.Println("Error reading response body: ", err)
				return
			}

			fmt.Println("Uploaded file: ", body)
		},
	}

	transfershCommand.Flags().StringP("url", "t", transfershUrl, "The URL to the transfer.sh instance")
	transfershCommand.Flags().StringP("username", "u", transfershUser, "The username to use for transfer.sh")
	transfershCommand.Flags().StringP("password", "p", transfershPass, "The password to use for transfer.sh")
	transfershCommand.Flags().StringP("filename", "f", "", "The filename to give the uploaded file")

	cmd.AddCommand(transfershCommand)
}
