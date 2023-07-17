package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var clusterAboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Cluster about",
	Long:  "Cluster about",
	RunE: func(cmd *cobra.Command, args []string) error {
		asRaw, _ := cmd.Flags().GetBool("raw")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		about, err := client.Cluster()
		if err != nil {
			return err
		}

		if asRaw {
			err := writeJSON(os.Stdout, about, true)
			return err
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"ID", "Name", "Version", "Uptime", "Last Contact", "Status", "CPU", "Memory", "Throttling"})

		for _, n := range about.Nodes {
			status := "follower"
			if n.Leader {
				status = "leader"
			}

			cpuusage := 0.0
			if n.Resources.CPULimit != 0 {
				cpuusage = n.Resources.CPU / n.Resources.CPULimit * 100
			}

			memoryusage := 0.0
			if n.Resources.MemLimit != 0 {
				memoryusage = float64(n.Resources.Mem) / float64(n.Resources.MemLimit) * 100
			}

			colors := text.Colors{text.Reset}
			if n.Uptime == 0 {
				colors = text.Colors{text.BgRed, text.BlinkSlow}
			}

			t.AppendRow(table.Row{
				colors.Sprint(n.ID),
				n.Name,
				n.Version,
				(time.Duration(n.Uptime) * time.Second).String(),
				(time.Duration(n.LastContact) * time.Millisecond).String(),
				status,
				fmt.Sprintf("%.1f%%", cpuusage),
				fmt.Sprintf("%.1f%%", memoryusage),
				fmt.Sprintf("%v", n.Resources.IsThrottling),
			})
		}

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 3, Align: text.AlignRight},
			{Number: 4, Align: text.AlignRight},
			{Number: 5, Align: text.AlignRight},
			{Number: 6, Align: text.AlignRight},
			{Number: 7, Align: text.AlignRight},
			{Number: 8, Align: text.AlignRight},
		})

		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Asc},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		if about.Degraded {
			fmt.Println(text.Colors{text.FgWhite, text.BgRed}.Sprint(about.DegradedErr))
		}

		fmt.Println("Raft:")

		t = table.NewWriter()

		t.AppendHeader(table.Row{"Address", "Last Contact", "Log Index", "Log Term", "Peers", "State"})

		colors := text.Colors{text.Reset}
		if about.Raft.LastContact > 200 || about.Raft.State == "Candidate" {
			colors = text.Colors{text.BgRed, text.BlinkSlow}
		}

		t.AppendRow(table.Row{
			about.Raft.Address,
			colors.Sprint((time.Duration(about.Raft.LastContact) * time.Millisecond).String()),
			about.Raft.LogIndex,
			about.Raft.LogTerm,
			about.Raft.NumPeers,
			colors.Sprint(about.Raft.State),
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil
	},
}

func init() {
	clusterCmd.AddCommand(clusterAboutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clusterAboutCmd.Flags().Bool("raw", false, "Display raw result from the API as JSON")
}
