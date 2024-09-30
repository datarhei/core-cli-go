package cmd

import (
	"fmt"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/spf13/cobra"
)

var clusterNodeFreeCmd = &cobra.Command{
	Use:   "free [id]",
	Short: "Move all processes from the node with the id to different nodes",
	Long:  "Move all processes from the node with the id to different nodes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nodeid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ClusterNodeProcessList(nodeid, coreclient.ProcessListOptions{})
		if err != nil {
			return err
		}

		total := len(list)
		for i, p := range list {
			id := coreclient.NewProcessID(p.ID, p.Domain)

			fmt.Printf("%4d / %4d: relocating %s away from %s ... ", i+1, total, id, nodeid)

			if err := client.ClusterRelocateProcess(id, ""); err != nil {
				fmt.Printf("failed: %s\n", err.Error())
			} else {
				fmt.Printf("OK\n")
			}
		}

		return nil
	},
}

func init() {
	clusterNodeCmd.AddCommand(clusterNodeFreeCmd)
}
