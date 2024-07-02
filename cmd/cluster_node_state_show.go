package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterNodeStateShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show the state of the node with the given id",
	Long:  "Show the state of the node with the given id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		state, err := client.ClusterNodeState(id)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, state.State, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeStateCmd.AddCommand(clusterNodeStateShowCmd)
}
