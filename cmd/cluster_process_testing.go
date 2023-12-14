package cmd

import (
	"fmt"
	"strconv"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

// processTestCmd represents the process command
var clusterProcessTestCmd = &cobra.Command{
	Use:   "test [number of processes] [owner] [source]",
	Short: "Process test",
	Long:  "Process test",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		owner := args[1]
		source := args[2]

		update, _ := cmd.Flags().GetBool("update")
		thumbs, _ := cmd.Flags().GetBool("thumbs")

		if n < 0 {
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

		fmt.Printf("%4d / %4d created\r", 0, n)

		for i := 0; i < n; i++ {
			name := "test_" + strconv.Itoa(i)

			config := api.ProcessConfig{
				ID:        name + "_main",
				Owner:     owner,
				Domain:    "",
				Type:      "ffmpeg",
				Reference: name,
				Input: []api.ProcessConfigIO{
					{
						ID:      "in_video",
						Address: source,
						Options: []string{
							"-thread_queue_size",
							"1024",
							"-re",
							"-copyts",
							"-start_at_zero",
							"-fflags",
							"+genpts+igndts",
						},
						Cleanup: []api.ProcessConfigIOCleanup{},
					},
					{
						ID:      "in_audio",
						Address: "anullsrc=r=44100:cl=mono",
						Options: []string{
							"-f",
							"lavfi",
							"-thread_queue_size",
							"1024",
							"-re",
						},
						Cleanup: []api.ProcessConfigIOCleanup{},
					},
				},
				Output: []api.ProcessConfigIO{
					{
						ID:      "out",
						Address: "{fs:mem}/" + name + "_%v.m3u8",
						Options: []string{
							"-c:v:0",
							"copy",
							"-bsf:v:0",
							"h264_metadata",
							"-c:a:0",
							"aac",
							"-map",
							"0:v:0",
							"-map",
							"1:a:0",
							"-f",
							"hls",
							"-start_number",
							"0",
							"-hls_time",
							"2",
							"-hls_list_size",
							"6",
							"-hls_delete_threshold",
							"12",
							"-hls_flags",
							"append_list+delete_segments+program_date_time+independent_segments",
							"-hls_segment_type",
							"mpegts",
							"-hls_segment_filename",
							"{fs:mem}/" + name + "_%v_%0004d.ts",
							"-master_pl_name",
							name + ".m3u8",
							"-master_pl_publish_rate",
							"5",
							"-var_stream_map",
							"v:0,a:0",
							"-y",
							"-method",
							"PUT",
							"-http_persistent",
							"1",
							"-ignore_io_errors",
							"1",
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

			if _, ok := processes[config.ID]; !ok {
				if !update {
					if err := client.ClusterProcessAdd(config); err != nil {
						fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", config.ID, i+1, n, err.Error())
						continue
					}
				}
			} else {
				if update {
					config.LogPatterns = append(config.LogPatterns, StringAlphanumeric(28))

					if err := client.ClusterProcessUpdate(coreclient.NewProcessID(config.ID, config.Domain), config); err != nil {
						fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", config.ID, i+1, n, err.Error())
					}
				}

				delete(processes, config.ID)
			}

			if thumbs {
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

				if _, ok := processes[config.ID]; !ok {
					if !update {
						if err := client.ClusterProcessAdd(config); err != nil {
							fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", config.ID, i+1, n, err.Error())
							continue
						}
					}
				} else {
					if update {
						config.LogPatterns = append(config.LogPatterns, StringAlphanumeric(28))

						if err := client.ClusterProcessUpdate(coreclient.NewProcessID(config.ID, config.Domain), config); err != nil {
							fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", config.ID, i+1, n, err.Error())
						}
					}

					delete(processes, config.ID)
				}
			}

			fmt.Printf("%4d / %4d created\r", i+1, n)
		}

		fmt.Printf("%4d / %4d created\n", n, n)

		i := 0
		n = len(processes)

		fmt.Printf("%4d / %4d deleted\r", 0, n)

		for name := range processes {
			if err := client.ClusterProcessDelete(coreclient.NewProcessID(name, "")); err != nil {
				fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", name, i+1, n, err.Error())
			}

			fmt.Printf("%4d / %4d deleted\r", i+1, n)
			i++
		}

		fmt.Printf("%4d / %4d deleted\n", n, n)

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessTestCmd)

	clusterProcessTestCmd.Flags().BoolP("update", "u", false, "Update existing processes")
	clusterProcessTestCmd.Flags().BoolP("thumbs", "t", false, "include thumbnail processes")
}
