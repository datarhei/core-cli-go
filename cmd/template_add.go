package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a template",
	Long:  "Add a template.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		from, _ := cmd.Flags().GetString("from")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		templates := viper.GetStringMap("templates")

		config := api.ProcessConfig{}

		if len(from) == 0 || from == "empty" {
			config = api.ProcessConfig{
				ID:        "",
				Owner:     "",
				Domain:    "",
				Type:      "ffmpeg",
				Reference: "",
				Input: []api.ProcessConfigIO{
					{
						ID:      "",
						Address: "",
						Options: []string{},
					},
				},
				Output: []api.ProcessConfigIO{
					{
						ID:      "",
						Address: "",
						Options: []string{},
						Cleanup: []api.ProcessConfigIOCleanup{
							{
								Pattern:       "",
								MaxFiles:      0,
								MaxFileAge:    0,
								PurgeOnDelete: false,
							},
						},
					},
				},
				Options:        []string{},
				Reconnect:      false,
				ReconnectDelay: 0,
				Autostart:      false,
				StaleTimeout:   0,
				Timeout:        0,
				LogPatterns:    []string{},
				Limits:         api.ProcessConfigLimits{},
				Metadata:       map[string]interface{}{},
			}
		} else if strings.HasPrefix(from, "process:") {
			client, err := connectSelectedCore()
			if err != nil {
				return err
			}

			processid := coreclient.ParseProcessID(strings.TrimPrefix(from, "process:"))

			process, err := client.Process(processid, []string{"config"})
			if err != nil {
				return err
			}

			config = *process.Config
		} else if strings.HasPrefix(from, "file:") {
			fromFile := strings.TrimPrefix(from, "file:")

			reader := os.Stdin

			if fromFile != "-" {
				file, err := os.Open(fromFile)
				if err != nil {
					return err
				}

				defer file.Close()

				reader = file
			}

			data, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &config)
			if err != nil {
				return err
			}
		}

		if !overwrite {
			if _, hasTemplate := templates[name]; hasTemplate {
				return fmt.Errorf("template with name '%s' already exists", name)
			}
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
	templateCmd.AddCommand(templateAddCmd)

	templateAddCmd.Flags().StringP("from", "f", "", "Source of template. 'empty', 'process:[id]', 'file:[path]'")
	templateAddCmd.Flags().BoolP("overwrite", "o", false, "Overwrite template")
}
