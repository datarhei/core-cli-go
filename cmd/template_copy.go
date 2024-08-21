package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateCopyCmd = &cobra.Command{
	Use:   "copy [from] [to]",
	Short: "Copy a template",
	Long:  "Copy a template.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		from := args[0]
		to := args[1]
		overwrite, _ := cmd.Flags().GetBool("overwrite")

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

		viper.Set("templates", templates)
		viper.WriteConfig()

		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateCopyCmd)

	templateCopyCmd.Flags().BoolP("overwrite", "o", false, "Overwrite template")
}
