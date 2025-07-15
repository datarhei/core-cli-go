package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterIamUserListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  "List all users in the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ClusterIdentitiesList(domain)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, list, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterIamUserCmd.AddCommand(clusterIamUserListCmd)

	clusterIamUserListCmd.Flags().StringP("domain", "d", "", "Domain")
}
