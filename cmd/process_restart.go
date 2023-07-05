package cmd

import (
	"fmt"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

// processRestartCmd represents the show command
var processRestartCmd = &cobra.Command{
	Use:   "restart [processid]",
	Short: "Restart the process with the given ID",
	Long:  "Restart the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		if err := client.ProcessCommand(id, "restart"); err != nil {
			return err
		}

		fmt.Printf("%s restart\n", id)

		return nil
	},
}

func init() {
	processCmd.AddCommand(processRestartCmd)
}
