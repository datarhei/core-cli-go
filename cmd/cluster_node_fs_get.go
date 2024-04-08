package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterNodeFilesystemGetCmd = &cobra.Command{
	Use:   "get [nodeid] [fsname] [path] [(-t|--to-file) path]",
	Short: "Download a file",
	Long:  "Download a file with the given path from the filesystem.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		storage := args[1]
		path := args[2]
		target, _ := cmd.Flags().GetString("to-file")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		file, err := client.ClusterNodeFilesystemGetFile(id, storage, path)
		if err != nil {
			return err
		}

		t := os.Stdout

		if target != "-" {
			file, err := os.Create(target)
			if err != nil {
				return err
			}

			t = file
			defer t.Close()
		}

		defer file.Close()

		t.ReadFrom(file)

		return nil
	},
}

func init() {
	clusterNodeFilesystemCmd.AddCommand(clusterNodeFilesystemGetCmd)

	clusterNodeFilesystemGetCmd.Flags().StringP("to-file", "t", "-", "Where to write the file, '-' for stdout")
}
