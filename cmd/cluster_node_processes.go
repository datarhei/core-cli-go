package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterNodeProcessesCmd = &cobra.Command{
	Use:   "processes [id]",
	Short: "Show the processes on the node with the given id",
	Long:  "Show the processes on the node with the given id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		processes, err := client.ClusterNodeProcessList(id)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, processes, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeCmd.AddCommand(clusterNodeProcessesCmd)
}
