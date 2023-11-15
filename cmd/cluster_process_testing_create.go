package cmd

import (
	"fmt"
	"strconv"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

var clusterProcessTestCreateCmd = &cobra.Command{
	Use:   "create [number of processes] [owner] [source]",
	Short: "Create processes",
	Long:  "Create processes.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		owner := args[1]
		source := args[2]

		if n <= 0 {
			return fmt.Errorf("the number of process must be positive")
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ClusterProcessList(coreclient.ProcessListOptions{
			IDPattern:     "test_*",
			RefPattern:    "test_*",
			OwnerPattern:  owner,
			DomainPattern: "",
		})
		if err != nil {
			return err
		}

		processes := map[string]struct{}{}

		for _, p := range list {
			processes[p.ID] = struct{}{}
		}

		fmt.Printf("%4d / %4d done\r", 0, n)

		for i := 0; i < n; i++ {
			name := "test_" + strconv.Itoa(i)

			if _, ok := processes[name+"_main"]; !ok {
				config := api.ProcessConfig{
					ID:        name + "_main",
					Owner:     owner,
					Domain:    "",
					Type:      "ffmpeg",
					Reference: name,
					Input: []api.ProcessConfigIO{
						{
							ID:      "in",
							Address: source,
							Options: []string{"-re"},
							Cleanup: []api.ProcessConfigIOCleanup{},
						},
					},
					Output: []api.ProcessConfigIO{
						{
							ID:      "out",
							Address: "{fs:mem}/" + name + ".m3u8",
							Options: []string{
								"-codec",
								"copy",
								"-f",
								"hls",
								"-start_number",
								"0",
								"-hls_time",
								"2",
								"-hls_list_size",
								"6",
								"-hls_flags",
								"append_list+delete_segments+program_date_time+temp_file",
								"-hls_delete_threshold",
								"4",
								"-hls_segment_filename",
								"{fs:mem}/" + name + "_%0004d.ts",
								"-y",
								"-method",
								"PUT",
							},
							Cleanup: []api.ProcessConfigIOCleanup{
								{
									Pattern:       "mem:/" + name + "_*",
									MaxFiles:      20,
									MaxFileAge:    0,
									PurgeOnDelete: true,
								},
							},
						},
					},
					Options:        []string{},
					Reconnect:      true,
					ReconnectDelay: 5,
					Autostart:      true,
					StaleTimeout:   10,
					Timeout:        0,
					Scheduler:      "",
					LogPatterns:    []string{},
					Limits: api.ProcessConfigLimits{
						CPU:     10,
						Memory:  50,
						WaitFor: 10,
					},
					Metadata: map[string]interface{}{},
				}

				if err := client.ClusterProcessAdd(config); err != nil {
					fmt.Printf("\nprocess %s_main (%4d / %4d) failed: %s\n", name, i+1, n, err.Error())
					continue
				}
			}
			/*
				if _, ok := processes[name+"_thumb"]; !ok {
					config := api.ProcessConfig{
						ID:        name + "_thumb",
						Owner:     owner,
						Domain:    "",
						Type:      "ffmpeg",
						Reference: name,
						Input: []api.ProcessConfigIO{
							{
								ID:      "in",
								Address: "{fs:mem}/" + name + ".m3u8",
								Options: []string{"-re"},
								Cleanup: []api.ProcessConfigIOCleanup{},
							},
						},
						Output: []api.ProcessConfigIO{
							{
								ID:      "jpeg",
								Address: "{fs:mem}/" + name + ".jpg",
								Options: []string{
									"-vframes", "1", "-method", "PUT", "-update", "1",
								},
							},
							{
								ID:      "jpeg_720",
								Address: "{fs:mem}/" + name + "_720.jpg",
								Options: []string{
									"-vframes", "1", "-vf",
									"scale=-1:720", "-method", "PUT", "-update", "1",
								},
							},
							{
								ID:      "jpeg_480",
								Address: "{fs:mem}/" + name + "_480.jpg",
								Options: []string{
									"-vframes", "1", "-vf",
									"scale=-1:480", "-method", "PUT", "-update", "1",
								},
							},
							{
								ID:      "jpeg_90",
								Address: "{fs:mem}/" + name + "_90.jpg",
								Options: []string{
									"-vframes", "1", "-vf",
									"scale=-1:90", "-method", "PUT", "-update", "1",
								},
							},
						},
						Options:        []string{},
						Reconnect:      true,
						ReconnectDelay: 60,
						Autostart:      true,
						StaleTimeout:   30,
						Timeout:        0,
						Scheduler:      "",
						LogPatterns:    []string{},
						Limits: api.ProcessConfigLimits{
							CPU:     10,
							Memory:  50,
							WaitFor: 10,
						},
						Metadata: map[string]interface{}{},
					}

					if err := client.ClusterProcessAdd(config); err != nil {
						fmt.Printf("\nprocess %s_thumb (%4d / %4d) failed: %s\n", name, i+1, n, err.Error())
					}
				}
			*/
			fmt.Printf("%4d / %4d done\r", i+1, n)
		}

		fmt.Printf("%4d / %4d done\n", n, n)

		return nil
	},
}

func init() {
	clusterProcessTestCmd.AddCommand(clusterProcessTestCreateCmd)

	//processAddCmd.Flags().String("from-file", "-", "Load process config from file or stdin")
}
