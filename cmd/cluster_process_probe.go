package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessProbeCmd = &cobra.Command{
	Use:   "probe [processid]",
	Short: "Probe the process with the given ID",
	Long:  "Probe the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		probe, err := client.ClusterProcessProbe(id)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, probe, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessProbeCmd)
}
