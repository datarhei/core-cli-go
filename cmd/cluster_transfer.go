package cmd

import (
	"github.com/spf13/cobra"
)

var clusterTransferCmd = &cobra.Command{
	Use:   "transfer [nodeid]",
	Short: "Transfer leadership to another node",
	Long:  "Transfer leadership to another node",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		return client.ClusterTransferLeadership(id)
	},
}

func init() {
	clusterCmd.AddCommand(clusterTransferCmd)
}
