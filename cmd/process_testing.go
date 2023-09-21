package cmd

import (
	"github.com/spf13/cobra"
)

// processTestCmd represents the process command
var processTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Process test related commands",
	Long:  "Process test related commands",
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("process called")
	//},
}

func init() {
	processCmd.AddCommand(processTestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//processCmd.PersistentFlags().Bool("raw", false, "Display raw result from the API as JSON")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
