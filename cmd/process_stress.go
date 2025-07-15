package cmd

import (
	"github.com/spf13/cobra"
)

var processStressCmd = &cobra.Command{
	Use:   "stress",
	Short: "Stress related commands",
	Long:  "Stress related commands",
}

func init() {
	processCmd.AddCommand(processStressCmd)
}
