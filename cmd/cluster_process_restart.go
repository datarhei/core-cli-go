package cmd

import (
	"fmt"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessRestartCmd = &cobra.Command{
	Use:   "restart [processid]",
	Short: "Restart the process with the given ID",
	Long:  "Restart the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		if err := client.ClusterProcessCommand(id, "restart"); err != nil {
			return err
		}

		fmt.Printf("%s restart\n", id)

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessRestartCmd)
}
