package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterDbPoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List all IAM policies",
	Long:  "List all IAM policies in the cluster DB",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ClusterDBPolicies()
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
	clusterDbCmd.AddCommand(clusterDbPoliciesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
