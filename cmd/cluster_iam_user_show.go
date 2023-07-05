package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterIamUserShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show the user with the given name",
	Long:  "Show the user with the given name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		policies, _ := cmd.Flags().GetBool("policies")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		user, err := client.ClusterIdentity(name)
		if err != nil {
			return err
		}

		if !policies {
			if err := writeJSON(os.Stdout, user, true); err != nil {
				return err
			}
		} else {
			if err := writeJSON(os.Stdout, user.Policies, true); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	clusterIamUserCmd.AddCommand(clusterIamUserShowCmd)

	clusterIamUserShowCmd.Flags().BoolP("policies", "p", false, "Show only the policies")
}
