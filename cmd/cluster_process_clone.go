package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var clusterProcessCloneCmd = &cobra.Command{
	Use:   "clone [processid]",
	Short: "Clone process config",
	Long:  "Clone the config of a process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		var fromClient coreclient.RestClient

		from, _ := cmd.Flags().GetString("from")
		if len(from) != 0 {
			otherClient, err := connectCore(from)
			if err != nil {
				return fmt.Errorf("connecting to %s: %w", from, err)
			}

			fromClient = otherClient
		} else {
			fromClient = client
		}

		id := coreclient.ParseProcessID(pid)

		process, err := fromClient.ClusterProcess(id, []string{"config"})
		if err != nil {
			return err
		}

		process.Config.ID += "_clone"

		data, err := json.MarshalIndent(process.Config, "", "   ")
		if err != nil {
			return err
		}

		editedData, modified, err := editData(data)
		if err != nil {
			return err
		}

		if !modified {
			fmt.Printf("No changes. Process config will not be cloned.\n")
			return nil
		}

		config := api.ProcessConfig{}

		if err := json.Unmarshal(editedData, &config); err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, config, true); err != nil {
			return err
		}

		return client.ClusterProcessAdd(config)
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessCloneCmd)

	clusterProcessCloneCmd.Flags().String("from", "", "Name of core to clone the process from")
}
