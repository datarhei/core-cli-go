package cmd

import (
	"fmt"
	"strconv"

	coreclient "github.com/datarhei/core-client-go/v16"
	coreclientapi "github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

// processTestCmd represents the process command
var processTestCmd = &cobra.Command{
	Use:   "test [template] [number of processes] [owner] [domain]?",
	Short: "Process test",
	Long:  "Process test",
	Args:  cobra.RangeArgs(3, 4),
	RunE: func(cmd *cobra.Command, args []string) error {
		template := args[0]
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		owner := args[2]
		domain := ""
		if len(args) == 4 {
			domain = args[3]
		}

		update, _ := cmd.Flags().GetBool("update")
		thumbs, _ := cmd.Flags().GetString("thumbs")
		autostart, _ := cmd.Flags().GetBool("autostart")
		metadata, _ := cmd.Flags().GetInt("metadata")

		if n < 0 {
			return fmt.Errorf("the number of process must be positive")
		}

		config, err := loadTemplate(template)
		if err != nil {
			return err
		}

		tconfig := coreclientapi.ProcessConfig{}

		if len(thumbs) != 0 {
			tconfig, err = loadTemplate(thumbs)
			if err != nil {
				return err
			}
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ProcessList(coreclient.ProcessListOptions{
			Filter:        []string{"state"},
			IDPattern:     "test_*",
			OwnerPattern:  owner,
			DomainPattern: domain,
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

			config.ID = name
			config.Owner = owner
			config.Domain = domain
			config.Autostart = autostart

			if metadata >= 0 {
				config.Metadata = map[string]interface{}{}
				config.Metadata["foobar"] = StringAlphanumeric(metadata * 1024)
			}

			if _, ok := processes[config.ID]; !ok {
				if !update {
					if err := client.ProcessAdd(config); err != nil {
						fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", config.ID, i+1, n, err.Error())
						continue
					}
				}
			} else {
				if update {
					config.LogPatterns = append(config.LogPatterns, StringAlphanumeric(28))

					if err := client.ProcessUpdate(coreclient.NewProcessID(config.ID, config.Domain), config); err != nil {
						fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", config.ID, i+1, n, err.Error())
					}
				}

				delete(processes, config.ID)
			}

			if len(thumbs) != 0 {
				tconfig.ID = name + "_thumb"
				tconfig.Owner = owner
				tconfig.Domain = domain
				tconfig.Autostart = true
				tconfig.Input[0].Address = "{fs:mem}/" + name + ".m3u8"

				if _, ok := processes[config.ID]; !ok {
					if !update {
						if err := client.ProcessAdd(config); err != nil {
							fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", config.ID, i+1, n, err.Error())
							continue
						}
					}
				} else {
					if update {
						config.LogPatterns = append(config.LogPatterns, StringAlphanumeric(28))

						if err := client.ProcessUpdate(coreclient.NewProcessID(config.ID, config.Domain), config); err != nil {
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
			if err := client.ProcessDelete(coreclient.NewProcessID(name, domain)); err != nil {
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
	processCmd.AddCommand(processTestCmd)

	processTestCmd.Flags().BoolP("update", "u", false, "Update existing processes")
	processTestCmd.Flags().StringP("thumbs", "t", "", "template for thumbnail processes")
	processTestCmd.Flags().BoolP("autostart", "a", true, "autostart processes")
	processTestCmd.Flags().IntP("metadata", "m", 0, "metadata size")
}
