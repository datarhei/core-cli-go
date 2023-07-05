package cmd

import (
	"github.com/spf13/cobra"
)

var clusterIamReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload cluster IAM",
	Long:  "Reload cluster IAM",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		return client.ClusterIAMReload()
	},
}

func init() {
	clusterIamCmd.AddCommand(clusterIamReloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//processCmd.PersistentFlags().Bool("raw", false, "Display raw result from the API as JSON")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
