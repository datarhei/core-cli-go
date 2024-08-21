package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateEditCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit a template",
	Long:  "Edit a template.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		templates := viper.GetStringMap("templates")

		config, hasTemplate := templates[name]
		if !hasTemplate {
			return fmt.Errorf("template with name '%s' doesn't exist", name)
		}

		data, err := json.MarshalIndent(config, "", "   ")
		if err != nil {
			return err
		}

		editedData, _, err := editData(data)
		if err != nil {
			return err
		}

		editedConfig := api.ProcessConfig{}

		if err := json.Unmarshal(editedData, &editedConfig); err != nil {
			return err
		}

		templates[name] = editedConfig

		viper.Set("templates", templates)
		viper.WriteConfig()

		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateEditCmd)
}
