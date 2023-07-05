package cmd

import (
	"fmt"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessReloadCmd = &cobra.Command{
	Use:   "reload [processid]",
	Short: "Reload the process with the given ID",
	Long:  "Reload the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		if err := client.ClusterProcessCommand(id, "reload"); err != nil {
			return err
		}

		fmt.Printf("%s reload\n", id)

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessReloadCmd)
}
