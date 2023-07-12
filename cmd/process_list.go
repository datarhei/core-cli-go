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

// processListCmd represents the list command
var processListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all processes",
	Long:  "List all processes of the selected core",
	RunE: func(cmd *cobra.Command, args []string) error {
		asRaw, _ := cmd.Flags().GetBool("raw")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		ids, _ := cmd.Flags().GetString("ids")
		filter, _ := cmd.Flags().GetString("filter")
		domain, _ := cmd.Flags().GetString("domain")
		reference, _ := cmd.Flags().GetString("reference")
		idpattern, _ := cmd.Flags().GetString("idpattern")
		refpattern, _ := cmd.Flags().GetString("refpattern")
		ownerpattern, _ := cmd.Flags().GetString("ownerpattern")
		domainpattern, _ := cmd.Flags().GetString("domainpattern")

		list, err := client.ProcessList(coreclient.ProcessListOptions{
			ID:            strings.Split(ids, ","),
			Filter:        strings.Split(filter, ","),
			Domain:        domain,
			Reference:     reference,
			IDPattern:     idpattern,
			RefPattern:    refpattern,
			OwnerPattern:  ownerpattern,
			DomainPattern: domainpattern,
		})
		if err != nil {
			return err
		}

		if asRaw {
			if err := writeJSON(os.Stdout, list, true); err != nil {
				return err
			}

			return nil
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"ID", "Domain", "Reference", "Order", "State", "Memory", "CPU", "Runtime", "Node", "Last Log"})

		for _, p := range list {
			runtime := p.State.Runtime
			if p.State.State != "running" {
				runtime = 0

				if p.State.Reconnect > 0 {
					runtime = -p.State.Reconnect
				}
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

			nodeid := ""
			if about, err := client.About(true); err == nil {
				nodeid = about.ID
			}

			lastlog := p.State.LastLog
			if len(lastlog) > 58 {
				lastlog = lastlog[:55] + "..."
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
				nodeid,
				lastlog,
			})
		}

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 2, Align: text.AlignRight},
			{Number: 4, Align: text.AlignRight},
			{Number: 5, Align: text.AlignRight},
			{Number: 6, Align: text.AlignRight},
			{Number: 7, Align: text.AlignRight},
			{Number: 8, Align: text.AlignRight},
		})

		t.SortBy([]table.SortBy{
			{Number: 2, Mode: table.Asc},
			{Number: 1, Mode: table.Asc},
			{Number: 4, Mode: table.Asc},
			{Number: 6, Mode: table.Dsc},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil
	},
}

func init() {
	processCmd.AddCommand(processListCmd)

	processListCmd.Flags().String("id", "", "A comma-separated list of process IDs")
	processListCmd.Flags().String("filter", "state", "A comma-separated list of filters per process: config, state, report, metadata")
	processListCmd.Flags().String("domain", "", "The domain to act upon")
	processListCmd.Flags().String("reference", "", "Limit list to specific reference")
	processListCmd.Flags().String("idpattern", "", "A glob pattern for the process IDs")
	processListCmd.Flags().String("refpattern", "", "A glob pattern for the process references")
	processListCmd.Flags().String("ownerpattern", "", "A gob pattern for the process owners")
	processListCmd.Flags().String("domainpattern", "", "A gob pattern for the process domains")
}
