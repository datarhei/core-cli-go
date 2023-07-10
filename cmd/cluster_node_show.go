package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterNodeShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show the node with the given id",
	Long:  "Show the node with the given id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		node, err := client.ClusterNode(id)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, node, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeCmd.AddCommand(clusterNodeShowCmd)
}
