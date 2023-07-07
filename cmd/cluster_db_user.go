package cmd

import (
	"github.com/spf13/cobra"
)

var clusterDbUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Cluster DB IAM user related commands",
	Long:  "Cluster DB IAM user related commands",
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("process called")
	//},
}

func init() {
	clusterDbCmd.AddCommand(clusterDbUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//processCmd.PersistentFlags().Bool("raw", false, "Display raw result from the API as JSON")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}