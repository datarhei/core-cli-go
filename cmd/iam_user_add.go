package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

// iamUserAddCmd represents the add command
var iamUserAddCmd = &cobra.Command{
	Use:   "add [name]?",
	Short: "Add a user",
	Long:  "Add a user to the core.",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := ""
		if len(args) == 1 {
			username = args[0]
		}

		var data []byte
		var err error

		if len(username) == 0 {
			fromFile, _ := cmd.Flags().GetString("from-file")
			if len(fromFile) == 0 {
				return fmt.Errorf("no user config file provided")
			}

			reader := os.Stdin

			if fromFile != "-" {
				file, err := os.Open(fromFile)
				if err != nil {
					return err
				}

				reader = file
			}

			data, err = io.ReadAll(reader)
			if err != nil {
				return err
			}
		} else {
			user := api.IAMUser{
				Name: username,
				Auth: api.IAMUserAuth{
					Services: api.IAMUserAuthServices{
						Basic:   []string{},
						Token:   []string{},
						Session: []string{},
					},
				},
				Policies: []api.IAMPolicy{
					{
						Domain:   "",
						Resource: "",
						Actions:  []string{},
					},
				},
			}

			data, err = json.MarshalIndent(user, "", "   ")
			if err != nil {
				return err
			}

			data, _, err = editData(data)
			if err != nil {
				return err
			}
		}

		config := api.IAMUser{}

		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		if err := client.IdentityAdd(config); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	iamUserCmd.AddCommand(iamUserAddCmd)

	iamUserAddCmd.Flags().String("from-file", "-", "Load user config from file or stdin")
}
