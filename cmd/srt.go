package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

// srtCmd represents the metrics command
var srtCmd = &cobra.Command{
	Use:   "srt",
	Short: "SRT related commands",
	Long:  "SRT related commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		asRaw, _ := cmd.Flags().GetBool("raw")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		var srt interface{}

		if asRaw {
			data, err := client.SRTChannelsRaw()
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &srt)
			if err != nil {
				return err
			}
		} else {
			srt, err = client.SRTChannels()
			if err != nil {
				return err
			}
		}

		if err := writeJSON(os.Stdout, srt, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(srtCmd)

	srtCmd.PersistentFlags().Bool("raw", false, "Display raw result from the API as JSON")
}
