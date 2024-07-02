package cmd

import (
	"github.com/spf13/cobra"
)

var clusterNodeStateSetCmd = &cobra.Command{
	Use:   "set [id] [state]",
	Short: "Set the state of the node with the given id",
	Long:  "Set the state of the node with the given id",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		state := args[1]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		err = client.ClusterNodeStateSet(id, state)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeStateCmd.AddCommand(clusterNodeStateSetCmd)
}
