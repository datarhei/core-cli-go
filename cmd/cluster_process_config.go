package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessConfigCmd = &cobra.Command{
	Use:   "config [processid]",
	Short: "Show the config of the process with the given ID",
	Long:  "Show the config of the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		process, err := client.ClusterProcess(id, []string{"config"})
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, process.Config, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessConfigCmd)
}
