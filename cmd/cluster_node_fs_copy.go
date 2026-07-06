package cmd

import (
	"github.com/spf13/cobra"
)

var clusterNodeFilesystemCopyCmd = &cobra.Command{
	Use:   "copy [nodeid] [srcfsname] [srcpath] [dstfsname] [dstpath]",
	Short: "Copy a file",
	Long:  "Copy a file to the same or different filesystem.",
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

		if err := client.ClusterNodeFilesystemCopyFile(id, dststorage, dstpath, srcstorage, srcpath); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeFilesystemCmd.AddCommand(clusterNodeFilesystemCopyCmd)
}
