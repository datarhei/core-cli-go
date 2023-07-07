package cmd

import (
	"fmt"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessDeleteCmd = &cobra.Command{
	Use:   "delete [processid] (-r|--reference)",
	Short: "Delete the process with the given ID",
	Long:  "Delete the process with the given ID or reference.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]
		reference, _ := cmd.Flags().GetBool("reference")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		if reference {
			list, err := client.ClusterProcessList(coreclient.ProcessListOptions{
				RefPattern: pid,
			})
			if err != nil {
				return err
			}

			for _, p := range list {
				id := coreclient.ProcessIDFromProcess(p)
				if err := client.ClusterProcessDelete(id); err != nil {
					fmt.Printf("%s error %s\n", id, err.Error())
				} else {
					fmt.Printf("%s delete\n", id)
				}
			}

			return nil
		}

		if err := client.ClusterProcessDelete(id); err != nil {
			return err
		}

		fmt.Printf("%s delete\n", id)

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessDeleteCmd)

	clusterProcessDeleteCmd.Flags().BoolP("reference", "r", false, "Interpret the processid as reference and delete all processes with that reference")
}