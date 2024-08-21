package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

var clusterProcessAddCmd = &cobra.Command{
	Use:   "add [processid]?",
	Short: "Add a process",
	Long:  "Add a process to the cluster.",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		processid := ""
		if len(args) == 1 {
			processid = args[0]
		}

		var data []byte
		var err error

		fromFile, _ := cmd.Flags().GetString("from-file")
		fromTemplate, _ := cmd.Flags().GetString("from-template")

		if len(fromTemplate) != 0 {
			config, err := loadTemplate(fromTemplate)
			if err != nil {
				return fmt.Errorf("loading from template: %w", err)
			}

			if len(processid) != 0 {
				config.ID = processid
			}

			data, err = json.MarshalIndent(config, "", "   ")
			if err != nil {
				return err
			}
		} else if len(fromFile) != 0 {
			reader := os.Stdin

			if fromFile != "-" {
				file, err := os.Open(fromFile)
				if err != nil {
					return err
				}

				reader = file
			}

			data, err = io.ReadAll(reader)
			if err != nil {
				return err
			}
		} else {
			config := api.ProcessConfig{
				ID:   processid,
				Type: "ffmpeg",
				Input: []api.ProcessConfigIO{
					{
						Options: []string{},
						Cleanup: []api.ProcessConfigIOCleanup{{}},
					},
				},
				Output: []api.ProcessConfigIO{
					{
						Options: []string{},
						Cleanup: []api.ProcessConfigIOCleanup{{}},
					},
				},
				Options:     []string{},
				LogPatterns: []string{},
				Metadata:    map[string]interface{}{},
			}

			data, err = json.MarshalIndent(config, "", "   ")
			if err != nil {
				return err
			}
		}

		data, _, err = editData(data)
		if err != nil {
			return err
		}

		config := api.ProcessConfig{}

		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		return client.ClusterProcessAdd(config)
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessAddCmd)

	clusterProcessAddCmd.Flags().String("from-file", "-", "Load process config from file or stdin")
	clusterProcessAddCmd.Flags().String("from-template", "", "Load process config from template")
}
