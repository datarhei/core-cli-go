package cmd

import (
	"os"

	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var clusterIamUserTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Print a template for a user config",
	Long:  "Print a template for a user config.",
	RunE: func(cmd *cobra.Command, args []string) error {
		user := api.IAMUser{
			Name: "TODO",
			Policies: []api.IAMPolicy{
				{
					Domain:   "$none",
					Resource: "",
					Actions:  []string{"any"},
				},
			},
		}

		if err := writeJSON(os.Stdout, user, false); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterIamUserCmd.AddCommand(clusterIamUserTemplateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
