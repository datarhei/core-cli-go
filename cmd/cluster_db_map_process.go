package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var clusterDbMapProcessCmd = &cobra.Command{
	Use:   "process",
	Short: "List a map of all processes and where they are currently deployed",
	Long:  "List a map of all processes and where they are currently deployed",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ClusterDBProcessMap()
		if err != nil {
			return err
		}

		invert, _ := cmd.Flags().GetBool("invert")
		if invert {
			m := map[string][]string{}

			for pid, nodeid := range list {
				a := append(m[nodeid], pid)
				m[nodeid] = a
			}

			return writeJSON(os.Stdout, m, true)
		} else {
			return writeJSON(os.Stdout, list, true)
		}
	},
}

func init() {
	clusterDbMapCmd.AddCommand(clusterDbMapProcessCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clusterDbMapProcessCmd.Flags().Bool("invert", false, "Invert the map")
}
