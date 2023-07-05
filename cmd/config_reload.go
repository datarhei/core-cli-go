package cmd

import (
	"github.com/spf13/cobra"
)

// configReloadCmd represents the list command
var configReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload the current config",
	Long:  "Reload the current config.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		if err := client.ConfigReload(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configReloadCmd)
}
