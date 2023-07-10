package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterNodeFilesCmd = &cobra.Command{
	Use:   "files [id]",
	Short: "Show the files on the node with the given id",
	Long:  "Show the files on the node with the given id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		files, err := client.ClusterNodeFiles(id)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, files, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeCmd.AddCommand(clusterNodeFilesCmd)
}
