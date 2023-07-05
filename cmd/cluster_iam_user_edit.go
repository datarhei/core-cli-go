package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var clusterIamUserEditCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit user config",
	Long:  "Edit user config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		policies, _ := cmd.Flags().GetBool("policies")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		user, err := client.ClusterIdentity(name)
		if err != nil {
			return err
		}

		var toEdit interface{}

		if !policies {
			toEdit = user
		} else {
			toEdit = user.Policies
		}

		data, err := json.MarshalIndent(toEdit, "", "   ")
		if err != nil {
			return err
		}

		editedData, modified, err := editData(data)
		if err != nil {
			return err
		}

		if !modified {
			// They are the same, nothing has been changed. No need to store the metadata
			fmt.Printf("No changes. User config will not be updated.")
			return nil
		}

		if !policies {
			config := api.IAMUser{}

			if err := json.Unmarshal(editedData, &config); err != nil {
				return err
			}

			if err := writeJSON(os.Stdout, config, true); err != nil {
				return err
			}

			return client.ClusterIdentityUpdate(name, config)
		} else {
			config := []api.IAMPolicy{}

			if err := json.Unmarshal(editedData, &config); err != nil {
				return err
			}

			if err := writeJSON(os.Stdout, config, true); err != nil {
				return err
			}

			return client.ClusterIdentitySetPolicies(name, config)
		}
	},
}

func init() {
	clusterIamUserCmd.AddCommand(clusterIamUserEditCmd)

	clusterIamUserEditCmd.Flags().BoolP("policies", "p", false, "Edit only the policies")
}
