package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clusterIamUserDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete the user with the given name",
	Long:  "Delete the user with the given name.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		if err := client.ClusterIdentityDelete(name); err != nil {
			return err
		}

		fmt.Printf("%s delete\n", name)

		return nil
	},
}

func init() {
	clusterIamUserCmd.AddCommand(clusterIamUserDeleteCmd)
}
