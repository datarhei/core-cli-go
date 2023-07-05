package cmd

import (
	"github.com/spf13/cobra"
)

var clusterLeaveCmd = &cobra.Command{
	Use:   "leave",
	Short: "Leave the cluster",
	Long:  "Leave the cluster",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		return client.ClusterLeave()
	},
}

func init() {
	clusterCmd.AddCommand(clusterLeaveCmd)
}
