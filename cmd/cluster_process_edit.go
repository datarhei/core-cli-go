package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var clusterProcessEditCmd = &cobra.Command{
	Use:   "edit [processid]",
	Short: "Edit process config",
	Long:  "Edit the config of a process",
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

		data, err := json.MarshalIndent(process.Config, "", "   ")
		if err != nil {
			return err
		}

		editedData, modified, err := editData(data)
		if err != nil {
			return err
		}

		if !modified {
			// They are the same, nothing has been changed. No need to store the metadata
			fmt.Printf("No changes. Process config will not be updated.")
			return nil
		}

		config := api.ProcessConfig{}

		if err := json.Unmarshal(editedData, &config); err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, config, true); err != nil {
			return err
		}

		return client.ClusterProcessUpdate(id, config)
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessEditCmd)
}
