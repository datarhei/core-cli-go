package cmd

import (
	"github.com/spf13/cobra"
)

var clusterLeaveCmd = &cobra.Command{
	Use:   "leave [id]",
	Short: "Leave the cluster",
	Long:  "Leave the cluster",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		return client.ClusterLeave(id)
	},
}

func init() {
	clusterCmd.AddCommand(clusterLeaveCmd)
}
