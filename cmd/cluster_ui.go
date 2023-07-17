package cmd

import (
	"fmt"

	"github.com/datarhei/core-cli-go/ui"

	"github.com/spf13/cobra"
)

var clusterUiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Cluster UI related commands",
	Long:  "Cluster UI related commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		_, err = client.Cluster()
		if err != nil {
			fmt.Println("cluster mode is not available")
			return err
		}

		return ui.Run(client)
	},
}

func init() {
	clusterCmd.AddCommand(clusterUiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//processCmd.PersistentFlags().Bool("raw", false, "Display raw result from the API as JSON")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
