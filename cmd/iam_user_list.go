package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// iamUserListCmd represents the list command
var iamUserListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  "List all users of the selected core",
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.IdentitiesList(domain)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, list, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	iamUserCmd.AddCommand(iamUserListCmd)

	iamUserListCmd.Flags().StringP("domain", "d", "", "Domain")
}
