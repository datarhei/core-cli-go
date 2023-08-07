package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:   "session [collector]",
	Short: "Session related commands",
	Long:  "Session related commands",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asRaw, _ := cmd.Flags().GetBool("raw")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.Sessions([]string{args[0]})
		if err != nil {
			return err
		}

		sessions := list[args[0]]

		if asRaw {
			if err := writeJSON(os.Stdout, sessions, true); err != nil {
				return err
			}

			return nil
		}

		data := map[string]session{}

		for _, sess := range sessions.Active.SessionList {
			s := data[sess.Location]

			s.count++
			s.bitrate += (sess.TxBitrate / 1024)

			data[sess.Location] = s
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Count", "Local", "Bitrate mbit"})

		sumBitrate := 0.0

		for l, sess := range data {
			t.AppendRow(table.Row{
				fmt.Sprintf("%5d", sess.count),
				l,
				sess.bitrate,
			})

			sumBitrate += sess.bitrate
		}

		t.AppendFooter(table.Row{
			fmt.Sprintf("%5d", len(sessions.Active.SessionList)),
			"",
			sumBitrate,
		})

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, Align: text.AlignRight},
			{
				Number: 3,
				Align:  text.AlignRight,
				Transformer: func(val interface{}) string {
					return fmt.Sprintf("%.3f", val)
				},
				TransformerFooter: func(val interface{}) string {
					return fmt.Sprintf("%.3f", val)
				},
			},
		})

		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Dsc},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	sessionCmd.PersistentFlags().Bool("raw", false, "Display raw result from the API as JSON")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
