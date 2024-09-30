package cmd

import (
	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessRelocateCmd = &cobra.Command{
	Use:   "relocate [processid] [targetnode]?",
	Short: "Move the process with the given ID to a different node",
	Long:  "Move the process with the given ID to a different node",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]
		target := ""
		if len(args) == 2 {
			target = args[1]
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		if err := client.ClusterRelocateProcess(id, target); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessRelocateCmd)
}
