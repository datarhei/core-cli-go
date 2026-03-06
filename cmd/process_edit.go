package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

// processEditCmd represents the list command
var processEditCmd = &cobra.Command{
	Use:   "edit [processid]",
	Short: "Edit process config",
	Long:  "Edit the config of a process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

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

		data, err := json.MarshalIndent(process.Config, "", "   ")
		if err != nil {
			return err
		}

		editedData, modified, err := editData(data)
		if err != nil {
			return err
		}

		if !force && !modified {
			// They are the same, nothing has been changed. No need to store the metadata
			fmt.Printf("No changes. Process config will not be updated.\n")
			return nil
		}

		config := api.ProcessConfig{}

		if err := json.Unmarshal(editedData, &config); err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, config, true); err != nil {
			return err
		}

		return client.ProcessUpdate(id, config, force)
	},
}

func init() {
	processCmd.AddCommand(processEditCmd)

	processEditCmd.Flags().Bool("force", false, "Whether to force an config update")
}
