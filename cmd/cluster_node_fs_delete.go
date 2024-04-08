package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clusterNodeFilesystemDeleteCmd = &cobra.Command{
	Use:   "delete [nodeid] [fsname] [path]",
	Short: "Delete a file",
	Long:  "Delete a file with the given path from the filesystem.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		storage := args[1]
		path := args[2]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		if err := client.ClusterNodeFilesystemDeleteFile(id, storage, path); err != nil {
			return err
		}

		fmt.Printf("%s:%s deleted\n", storage, path)

		return nil
	},
}

func init() {
	clusterNodeFilesystemCmd.AddCommand(clusterNodeFilesystemDeleteCmd)
}
