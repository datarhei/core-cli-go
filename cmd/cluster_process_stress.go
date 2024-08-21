package cmd

import (
	"github.com/spf13/cobra"
)

var clusterProcessStressCmd = &cobra.Command{
	Use:   "stress",
	Short: "Stress related commands",
	Long:  "Stress related commands",
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessStressCmd)
}
