package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fsDeleteCmd represents the list command
var fsDeleteCmd = &cobra.Command{
	Use:   "delete [fsname] [path]",
	Short: "Delete a file",
	Long:  "Delete a file with the given path from the filesystem.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := args[1]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		if err := client.FilesystemDeleteFile(name, path); err != nil {
			return err
		}

		fmt.Printf("%s:%s deleted\n", name, path)

		return nil
	},
}

func init() {
	fsCmd.AddCommand(fsDeleteCmd)
}
