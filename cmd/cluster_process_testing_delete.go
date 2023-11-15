package cmd

import (
	"fmt"
	"strconv"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/spf13/cobra"
)

var clusterProcessTestDeleteCmd = &cobra.Command{
	Use:   "delete [number of processes] [owner]",
	Short: "Delete processes",
	Long:  "Delete processes.",
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

		if n <= 0 {
			i := 0
			n = len(processes)

			fmt.Printf("%4d / %4d done\r", 0, n)

			for name := range processes {
				if err := client.ClusterProcessDelete(coreclient.NewProcessID(name, "")); err != nil {
					fmt.Printf("\nprocess %s (%4d / %4d) failed: %s\n", name, i+1, n, err.Error())
				}

				fmt.Printf("%4d / %4d done\r", i+1, n)
				i++
			}
		} else {
			fmt.Printf("%4d / %4d done\r", 0, n)

			for i := 0; i < n; i++ {
				name := "test_" + strconv.Itoa(i)

				if _, ok := processes[name+"_main"]; ok {
					if err := client.ClusterProcessDelete(coreclient.NewProcessID(name+"_main", "")); err != nil {
						fmt.Printf("\nprocess %s_main (%4d / %4d) failed: %s\n", name, i+1, n, err.Error())
					}
				}

				if _, ok := processes[name+"_thumb"]; ok {
					if err := client.ClusterProcessDelete(coreclient.NewProcessID(name+"_thumb", "")); err != nil {
						fmt.Printf("\nprocess %s_thumb (%4d / %4d) failed: %s\n", name, i+1, n, err.Error())
					}
				}

				fmt.Printf("%4d / %4d done\r", i+1, n)
			}
		}

		fmt.Printf("%4d / %4d done\n", n, n)

		return nil
	},
}

func init() {
	clusterProcessTestCmd.AddCommand(clusterProcessTestDeleteCmd)

	//processAddCmd.Flags().String("from-file", "-", "Load process config from file or stdin")
}
