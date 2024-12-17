package txt

import (
	"io"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/v3-nu/vvv/cmd/utils"
)

func ReplaceText(cmd *cobra.Command) {
	command := &cobra.Command{
		Use:     "replace",
		Aliases: []string{"repl"},
		Short:   "Replace text in stdin and print to stdout",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			regex := cmd.Flag("regex").Value.String() == "true"
			srcFile := cmd.Flag("src").Value.String()
			dstFile := cmd.Flag("dst").Value.String()

			var ch io.Reader
			var src io.Reader
			var dst io.Writer
			var err error

			if srcFile != "" {
				src, err = os.Open(srcFile)
				if err != nil {
					log.Fatalf("Failed to open source file: %v", err)
				}
			} else {
				src = os.Stdin
			}

			if dstFile != "" {
				dst, err = os.Create(dstFile)
				if err != nil {
					log.Fatalf("Failed to create destination file: %v", err)
				}
			} else {
				dst = os.Stdout
			}

			if regex {
				regexp := regexp.MustCompile(args[0])

				ch = Chain(src, Regexp(regexp, []byte(args[1])))
			} else {
				ch = Chain(src, String(args[0], args[1]))
			}

			io.Copy(dst, ch)
		},
	}

	command.Flags().BoolP("regex", "r", false, "Use regex replacement")
	command.Flags().StringP("src", "s", "", "Source file")
	command.Flags().StringP("dst", "d", "", "Destination file")

	cmd.AddCommand(command)
}

var ExportCommands = utils.CommandGroup{
	Command: &cobra.Command{
		Use:   "txt",
		Short: "Text related commands",
	},
	Children: []func(*cobra.Command){
		ReplaceText,
	},
}
