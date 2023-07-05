package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Download a snapshot of the cluster DB",
	Long:  "Download a snapshot of the cluster DB",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("to-file")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		file, err := client.ClusterSnapshot()
		if err != nil {
			return err
		}

		defer file.Close()

		t := os.Stdout

		if target != "-" {
			file, err := os.Create(target)
			if err != nil {
				return err
			}

			t = file
			defer t.Close()
		}

		t.ReadFrom(file)

		return nil
	},
}

func init() {
	clusterCmd.AddCommand(clusterSnapshotCmd)

	clusterSnapshotCmd.Flags().StringP("to-file", "t", "-", "Where to write the file, '-' for stdout")
}
