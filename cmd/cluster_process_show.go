package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var clusterProcessShowCmd = &cobra.Command{
	Use:   "show [processid]",
	Short: "Show the process with the given ID",
	Long:  "Show the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]
		asRaw, _ := cmd.Flags().GetBool("raw")
		withConfig, _ := cmd.Flags().GetBool("cfg")
		withReport, _ := cmd.Flags().GetBool("report")
		withMetadata, _ := cmd.Flags().GetBool("metadata")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		filter := []string{"state"}
		if withConfig {
			filter = append(filter, "config")
		}
		if withReport {
			filter = append(filter, "report")
		}
		if withMetadata {
			filter = append(filter, "metadata")
		}

		p, err := client.ClusterProcess(id, filter)
		if err != nil {
			return err
		}

		if asRaw {
			if err := writeJSON(os.Stdout, p, true); err != nil {
				return err
			}

			return nil
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"ID", "Domain", "Reference", "Order", "State", "Memory", "CPU", "Runtime"})

		runtime := p.State.Runtime
		if p.State.State != "running" {
			runtime = 0
		}

		order := strings.ToUpper(p.State.Order)
		switch order {
		case "START":
			order = text.FgGreen.Sprint(order)
		case "STOP":
			order = text.Colors{text.FgWhite, text.Faint}.Sprint(order)
		}

		state := strings.ToUpper(p.State.State)
		switch state {
		case "RUNNING":
			state = text.FgGreen.Sprint(state)
		case "FINISHED":
			state = text.Colors{text.FgWhite, text.Faint}.Sprint(state)
		case "FAILED":
			state = text.FgRed.Sprint(state)
		case "STARTING":
			state = text.FgCyan.Sprint(state)
		case "FINISHING":
			state = text.FgCyan.Sprint(state)
		case "KILLED":
			state = text.Colors{text.FgRed, text.Faint}.Sprint(state)
		}

		t.AppendRow(table.Row{
			p.ID,
			p.Domain,
			p.Reference,
			order,
			state,
			formatByteCountBinary(p.State.Memory),
			fmt.Sprintf("%.1f%%", p.State.CPU),
			(time.Duration(runtime) * time.Second).String(),
		})

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 2, Align: text.AlignRight},
			{Number: 4, Align: text.AlignRight},
			{Number: 5, Align: text.AlignRight},
			{Number: 6, Align: text.AlignRight},
			{Number: 7, Align: text.AlignRight},
			{Number: 8, Align: text.AlignRight},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		if p.State == nil || p.State.Progress == nil {
			return nil
		}

		if len(p.State.Progress.Input) == 0 && len(p.State.Progress.Output) == 0 {
			return nil
		}

		t = table.NewWriter()

		rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

		t.SetTitle("Inputs / Outputs")
		t.AppendHeader(table.Row{"", "#", "ID", "Type", "URL", "Specs"}, rowConfigAutoMerge)

		for i, p := range p.State.Progress.Input {
			var specs string
			if p.Type == "audio" {
				specs = fmt.Sprintf("%s %s %dHz", strings.ToUpper(p.Codec), p.Layout, p.Sampling)
			} else {
				specs = fmt.Sprintf("%s %dx%d", strings.ToUpper(p.Codec), p.Width, p.Height)
			}

			t.AppendRow(table.Row{
				"input",
				i,
				p.ID,
				strings.ToUpper(p.Type),
				p.Address,
				specs,
			}, rowConfigAutoMerge)
		}

		for i, p := range p.State.Progress.Output {
			var specs string
			if p.Type == "audio" {
				specs = fmt.Sprintf("%s %s %dHz", strings.ToUpper(p.Codec), p.Layout, p.Sampling)
			} else {
				specs = fmt.Sprintf("%s %dx%d", strings.ToUpper(p.Codec), p.Width, p.Height)
			}

			t.AppendRow(table.Row{
				"output",
				i,
				p.ID,
				strings.ToUpper(p.Type),
				p.Address,
				specs,
			}, rowConfigAutoMerge)
		}

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessShowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	clusterProcessShowCmd.Flags().BoolP("cfg", "c", false, "Include the process config")
	clusterProcessShowCmd.Flags().BoolP("report", "r", false, "Include the process config")
	clusterProcessShowCmd.Flags().BoolP("metadata", "m", false, "Include the process config")
}
