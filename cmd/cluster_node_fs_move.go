package cmd

import (
	"github.com/spf13/cobra"
)

var clusterNodeFilesystemMoveCmd = &cobra.Command{
	Use:   "move [nodeid] [srcfsname] [srcpath] [dstfsname] [dstpath]",
	Short: "Move a file",
	Long:  "Move a file to the same or different filesystem.",
	Args:  cobra.ExactArgs(5),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		srcstorage := args[1]
		srcpath := args[2]
		dststorage := args[3]
		dstpath := args[4]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		if err := client.ClusterNodeFilesystemMoveFile(id, dststorage, dstpath, srcstorage, srcpath); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeFilesystemCmd.AddCommand(clusterNodeFilesystemMoveCmd)
}
