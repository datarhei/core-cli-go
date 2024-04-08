package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterNodeFilesystemAddCmd = &cobra.Command{
	Use:   "add [nodeid] [fsname] [path] [(-f|--from-file) path]",
	Short: "Upload a file",
	Long:  "Upload a file with the given path from the filesystem.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		storage := args[1]
		path := args[2]
		source, _ := cmd.Flags().GetString("from-file")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		s := os.Stdin

		if source != "-" {
			file, err := os.Open(source)
			if err != nil {
				return err
			}

			s = file
			defer s.Close()
		}

		if err := client.ClusterNodeFilesystemPutFile(id, storage, path, s); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	clusterNodeFilesystemCmd.AddCommand(clusterNodeFilesystemAddCmd)

	clusterNodeFilesystemAddCmd.Flags().StringP("from-file", "f", "-", "Where to read the file from, '-' for stdin")
}
