package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterAboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Cluster about",
	Long:  "Cluster about",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		about, err := client.Cluster()
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, about, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterCmd.AddCommand(clusterAboutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
