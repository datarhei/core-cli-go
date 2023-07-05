package cmd

import (
	"fmt"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

// stooProcessCmd represents the show command
var processStopCmd = &cobra.Command{
	Use:   "stop [processid]",
	Short: "Stop the process with the given ID",
	Long:  "Stop the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		if err := client.ProcessCommand(id, "stop"); err != nil {
			return err
		}

		fmt.Printf("%s stop\n", id)

		return nil
	},
}

func init() {
	processCmd.AddCommand(processStopCmd)
}
