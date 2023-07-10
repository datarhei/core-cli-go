package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterNodeVersionCmd = &cobra.Command{
	Use:   "version [id]",
	Short: "Show the version of the node with the given id",
	Long:  "Show the version of the node with the given id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		version, err := client.ClusterNodeVersion(id)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, version, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeCmd.AddCommand(clusterNodeVersionCmd)
}
