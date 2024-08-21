package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateMoveCmd = &cobra.Command{
	Use:   "move [from] [to]",
	Short: "Rename a template",
	Long:  "Rename a template.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		from := args[0]
		to := args[1]
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		if to == from {
			return nil
		}

		templates := viper.GetStringMap("templates")

		fromConfig, hasTemplate := templates[from]
		if !hasTemplate {
			return fmt.Errorf("the template with name '%s' doesn't exist", from)
		}

		if !overwrite {
			if _, hasTemplate := templates[to]; hasTemplate {
				return fmt.Errorf("the template with the name '%s' already exists", to)
			}
		}

		templates[to] = fromConfig
		delete(templates, from)

		viper.Set("templates", templates)
		viper.WriteConfig()

		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateMoveCmd)

	templateMoveCmd.Flags().BoolP("overwrite", "o", false, "Overwrite template")
}
