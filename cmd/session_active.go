package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

type session struct {
	count      uint64
	rx_bitrate float64 // mbit/session
	tx_bitrate float64 // mbit/session
}

var sessionActiveCmd = &cobra.Command{
	Use:   "active [collector]",
	Short: "List all active sessions for a collector",
	Long:  "List all active sessions for a collector",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asRaw, _ := cmd.Flags().GetBool("raw")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.SessionsActive([]string{args[0]})
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

		for _, sess := range sessions {
			s := data[sess.Location]

			s.count++
			s.rx_bitrate += (sess.RxBitrate / 1024)
			s.tx_bitrate += (sess.TxBitrate / 1024)

			data[sess.Location] = s
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Count", "Local", "RX bitrate mbit", "TX bitrate mbit"})

		sumRxBitrate := 0.0
		sumTxBitrate := 0.0

		for l, sess := range data {
			t.AppendRow(table.Row{
				fmt.Sprintf("%5d", sess.count),
				l,
				sess.rx_bitrate,
				sess.tx_bitrate,
			})

			sumRxBitrate += sess.rx_bitrate
			sumTxBitrate += sess.tx_bitrate
		}

		t.AppendFooter(table.Row{
			fmt.Sprintf("%5d", len(sessions)),
			"",
			sumRxBitrate,
			sumTxBitrate,
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
	sessionCmd.AddCommand(sessionActiveCmd)
}
