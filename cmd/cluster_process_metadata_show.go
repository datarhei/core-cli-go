package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessMetadataShowCmd = &cobra.Command{
	Use:   "show [processid] [key]?",
	Short: "Show the metadata of the process with the given ID",
	Long:  "Show the metadata of the process with the given ID",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]
		key := ""
		if len(args) == 2 {
			key = args[1]
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		metadata, err := client.ClusterProcessMetadata(id, key)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, metadata, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterProcessMetadataCmd.AddCommand(clusterProcessMetadataShowCmd)
}
