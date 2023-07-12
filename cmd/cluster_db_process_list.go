package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterDbProcessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all processes",
	Long:  "List all processes in the cluster DB",
	RunE: func(cmd *cobra.Command, args []string) error {
		asRaw, _ := cmd.Flags().GetBool("raw")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ClusterDBProcessList()
		if err != nil {
			return err
		}

		if asRaw {
			if err := writeJSON(os.Stdout, list, true); err != nil {
				return err
			}

			return nil
		}

		pmap, err := client.ClusterDBProcessMap()
		if err != nil {
			return err
		}

		dbProcessTable(list, pmap)

		return nil
	},
}

func init() {
	clusterDbProcessCmd.AddCommand(clusterDbProcessListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
