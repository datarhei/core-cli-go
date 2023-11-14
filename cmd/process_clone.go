package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var processCloneCmd = &cobra.Command{
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

		id := coreclient.ParseProcessID(pid)

		process, err := client.Process(id, []string{"config"})
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
			fmt.Printf("No changes. Process config will not be cloned.")
			return nil
		}

		config := api.ProcessConfig{}

		if err := json.Unmarshal(editedData, &config); err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, config, true); err != nil {
			return err
		}

		return client.ProcessAdd(config)
	},
}

func init() {
	processCmd.AddCommand(processCloneCmd)
}
