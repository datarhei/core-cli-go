package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show a template",
	Long:  "Show a template.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		templates := viper.GetStringMap("templates")

		config, hasTemplate := templates[name]
		if !hasTemplate {
			return fmt.Errorf("template with name '%s' doesn't exist", name)
		}

		return writeJSON(os.Stdout, config, true)
	},
}

func init() {
	templateCmd.AddCommand(templateShowCmd)
}
