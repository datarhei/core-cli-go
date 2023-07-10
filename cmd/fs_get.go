package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// fsGetCmd represents the list command
var fsGetCmd = &cobra.Command{
	Use:   "get [fsname] [path] [(-t|--to-file) path]",
	Short: "Download a file",
	Long:  "Download a file with the given path from the filesystem.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := args[1]
		target, _ := cmd.Flags().GetString("to-file")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		file, err := client.FilesystemGetFile(name, path)
		if err != nil {
			return err
		}

		t := os.Stdout

		if target != "-" {
			file, err := os.Create(target)
			if err != nil {
				return err
			}

			t = file
			defer t.Close()
		}

		defer file.Close()

		t.ReadFrom(file)

		return nil
	},
}

func init() {
	fsCmd.AddCommand(fsGetCmd)

	fsGetCmd.Flags().StringP("to-file", "t", "-", "Where to write the file, '-' for stdout")
}
