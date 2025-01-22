package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var processPlayoutReopenCmd = &cobra.Command{
	Use:   "reopen [processid] [inputid]",
	Short: "Show playout reopen",
	Long:  "Show playout reopen",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]
		inputid := args[1]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		status, err := client.PlayoutReopen(id, inputid)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, status, true); err != nil {
			return err
		}

		return nil

	},
}

func init() {
	processPlayoutCmd.AddCommand(processPlayoutReopenCmd)
}
