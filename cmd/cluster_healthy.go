package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterHealthyCmd = &cobra.Command{
	Use:   "healthy",
	Short: "Healthiness of the cluster",
	Long:  "Healthiness of the cluster",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		healthy, err := client.ClusterHealthy()
		if err != nil {
			return err
		}

		err = writeJSON(os.Stdout, healthy, true)

		return err
	},
}

func init() {
	clusterCmd.AddCommand(clusterHealthyCmd)
}
