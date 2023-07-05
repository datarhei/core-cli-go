package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// fsAddCmd represents the list command
var fsAddCmd = &cobra.Command{
	Use:   "add [name] [path] [(-f|--from-file) path]",
	Short: "Upload a file",
	Long:  "Upload a file with the given path from the filesystem.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := args[1]
		source, _ := cmd.Flags().GetString("from-file")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		s := os.Stdin

		if source != "-" {
			file, err := os.Open(source)
			if err != nil {
				return err
			}

			s = file
			defer s.Close()
		}

		if err := client.FilesystemAddFile(name, path, s); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	fsCmd.AddCommand(fsAddCmd)

	fsAddCmd.Flags().StringP("from-file", "f", "-", "Where to read the file from, '-' for stdin")
}
