package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessStateCmd = &cobra.Command{
	Use:   "state [processid]",
	Short: "Show the state of the process with the given ID",
	Long:  "Show the state of the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		process, err := client.ClusterProcess(id, []string{"state"})
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, process.State, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessStateCmd)
}
