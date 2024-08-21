package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a template",
	Long:  "Delete a template.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		templates := viper.GetStringMap("templates")

		delete(templates, name)

		viper.Set("templates", templates)
		viper.WriteConfig()

		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateDeleteCmd)
}
