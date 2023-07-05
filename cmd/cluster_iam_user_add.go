package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

var clusterIamUserAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a user",
	Long:  "Add a user to the cluster.",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		data, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		config := api.IAMUser{}

		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		if err := client.ClusterIdentityAdd(config); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterIamUserCmd.AddCommand(clusterIamUserAddCmd)

	clusterIamUserAddCmd.Flags().String("from-file", "-", "Load user config from file or stdin")
}
