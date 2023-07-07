package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/spf13/cobra"
)

var clusterDbProcessShowCmd = &cobra.Command{
	Use:   "show [processid]",
	Short: "Show a specific process",
	Long:  "Show a specific process in the cluster DB",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		process, err := client.ClusterDBProcess(id)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, process, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterDbProcessCmd.AddCommand(clusterDbProcessShowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}