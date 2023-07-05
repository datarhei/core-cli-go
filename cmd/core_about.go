package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// coreAboutCmd represents the backup command
var coreAboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Show core details",
	Long:  "Show core details.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		about, err := client.About(true)
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, about, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	coreCmd.AddCommand(coreAboutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
