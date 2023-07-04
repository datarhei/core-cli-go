package cmd

import (
	"fmt"
	"strings"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

// stateProcessCmd represents the show command
var commandProcessCmd = &cobra.Command{
	Use:   "command [processid]",
	Short: "Show the ffmpeg command of the process with the given ID",
	Long:  "Show the ffmpeg command of the process with the given ID",
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

		for _, e := range state.Command {
			if strings.ContainsAny(e, " $") {
				fmt.Printf("'%s' ", e)
			} else {
				fmt.Printf("%s ", e)
			}
		}

		fmt.Printf("\n")

		return nil
	},
}

func init() {
	processCmd.AddCommand(commandProcessCmd)
}
