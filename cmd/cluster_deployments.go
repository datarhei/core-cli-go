package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterDeploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "Currently pending deployments",
	Long:  "Currently pending deployments",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		deployments, err := client.ClusterDeployments()
		if err != nil {
			return err
		}

		err = writeJSON(os.Stdout, deployments, true)

		return err
	},
}

func init() {
	clusterCmd.AddCommand(clusterDeploymentsCmd)
}
