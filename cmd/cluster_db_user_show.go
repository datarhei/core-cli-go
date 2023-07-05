package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterDbUserShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show a specific IAM user",
	Long:  "Show a specific IAM user in the cluster DB",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		user, err := client.ClusterDBUser(name)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, user, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterDbUserCmd.AddCommand(clusterDbUserShowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
