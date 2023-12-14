package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var iamUserCloneCmd = &cobra.Command{
	Use:   "clone [name]",
	Short: "Clone user config",
	Long:  "Clone user config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		user, err := client.Identity(name)
		if err != nil {
			return err
		}

		user.Name += "_clone"
		user.Alias += "_clone"

		data, err := json.MarshalIndent(user, "", "   ")
		if err != nil {
			return err
		}

		editedData, modified, err := editData(data)
		if err != nil {
			return err
		}

		if !modified {
			// They are the same, nothing has been changed
			fmt.Printf("No changes. User config will not be cloned.\n")
			return nil
		}

		config := api.IAMUser{}

		if err := json.Unmarshal(editedData, &config); err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, config, true); err != nil {
			return err
		}

		return client.IdentityAdd(config)
	},
}

func init() {
	iamUserCmd.AddCommand(iamUserCloneCmd)
}
