package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

// stateProcessCmd represents the show command
var stateProcessCmd = &cobra.Command{
	Use:   "state [processid]",
	Short: "Show the state of the process with the given ID",
	Long:  "Show the state of the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		state, err := client.ProcessState(id)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, state, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	processCmd.AddCommand(stateProcessCmd)
}
