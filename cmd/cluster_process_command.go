package cmd

import (
	"fmt"
	"strings"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

var clusterProcessCommandCmd = &cobra.Command{
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

		process, err := client.ClusterProcess(id, []string{"state"})
		if err != nil {
			return err
		}

		if len(process.State.Command) == 0 {
			return fmt.Errorf("the command is not available")
		}

		for _, e := range process.State.Command {
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
	clusterProcessCmd.AddCommand(clusterProcessCommandCmd)
}
