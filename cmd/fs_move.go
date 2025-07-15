package cmd

import (
	"github.com/spf13/cobra"
)

var fsMoveCmd = &cobra.Command{
	Use:   "move [srcfsname] [srcpath] [dstfsname] [dstpath]",
	Short: "Move/rename a file",
	Long:  "Move/rename a file to the same or different filesystem.",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcstorage := args[0]
		srcpath := args[1]
		dststorage := args[2]
		dstpath := args[3]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		if err := client.FilesystemMoveFile(dststorage, dstpath, srcstorage, srcpath); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	fsCmd.AddCommand(fsMoveCmd)
}
