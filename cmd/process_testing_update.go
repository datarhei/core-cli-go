package cmd

import (
	"fmt"
	"strconv"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

// processAddCmd represents the add command
var processTestUpdateCmd = &cobra.Command{
	Use:   "update [number of processes] [owner]",
	Short: "Update processes",
	Long:  "Update processes.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		owner := args[1]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			fmt.Printf("%4d / %4d done\r", i+1, n)

			config := api.ProcessConfig{
				ID:        "processTest-" + strconv.Itoa(i),
				Owner:     owner,
				Domain:    "",
				Type:      "ffmpeg",
				Reference: StringAlphanumeric(28),
				Input: []api.ProcessConfigIO{
					{
						ID:      "in",
						Address: "testsrc2=rate=25:size=640x360",
						Options: []string{
							"-f", "lavfi",
							"-re",
						},
						Cleanup: []api.ProcessConfigIOCleanup{},
					},
				},
				Output: []api.ProcessConfigIO{
					{
						ID:      "out",
						Address: "-",
						Options: []string{
							"-codec", "copy",
							"-f", "null",
						},
						Cleanup: []api.ProcessConfigIOCleanup{},
					},
				},
				Options:        []string{},
				Reconnect:      true,
				ReconnectDelay: 10,
				Autostart:      true,
				StaleTimeout:   10,
				Timeout:        0,
				Scheduler:      "",
				LogPatterns:    []string{},
				Limits: api.ProcessConfigLimits{
					CPU:     5,
					Memory:  50,
					WaitFor: 0,
				},
				Metadata: map[string]interface{}{},
			}

			if err := client.ClusterProcessUpdate(coreclient.NewProcessID(config.ID, config.Domain), config); err != nil {
				return err
			}
		}

		fmt.Printf("%4d / %4d done\r", n, n)

		return nil
	},
}

func init() {
	processTestCmd.AddCommand(processTestUpdateCmd)

	//processAddCmd.Flags().String("from-file", "-", "Load process config from file or stdin")
}
