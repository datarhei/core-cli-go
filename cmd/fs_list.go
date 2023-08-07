package cmd

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// fsListCmd represents the list command
var fsListCmd = &cobra.Command{
	Use:   "list [fsname] [pattern]? (-s|--sort) [none|name|size|lastmod] (-o|--order) [asc|desc] (-t|--target) [url with %%s]",
	Short: "List files",
	Long:  "List files on filesystem",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		pattern := ""
		if len(args) == 2 {
			pattern = args[1]
		}

		sort, _ := cmd.Flags().GetString("sort")
		order, _ := cmd.Flags().GetString("order")
		target, _ := cmd.Flags().GetString("target")
		random, _ := cmd.Flags().GetBool("random")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.FilesystemList(name, pattern, sort, order)
		if err != nil {
			return err
		}

		if len(target) != 0 {
			for _, f := range list {
				if random {
					fmt.Printf(target+"\n", "/"+StringAlphanumeric(12))
				} else {
					fmt.Printf(target+"\n", f.Name)
				}
			}

			return nil
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Name", "Size", "Last Modification"})

		for _, f := range list {
			lastMod := time.Unix(f.LastMod, 0)
			t.AppendRow(table.Row{f.Name, formatByteCountBinary(uint64(f.Size)), lastMod.Format("2006-01-02 15:04:05")})
		}

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 2, Align: text.AlignRight},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil

	},
}

func init() {
	fsCmd.AddCommand(fsListCmd)

	fsListCmd.Flags().StringP("sort", "s", "none", "Sorting criteria")
	fsListCmd.Flags().StringP("order", "o", "asc", "Sorting direction")
	fsListCmd.Flags().StringP("target", "t", "", "Create vegeta targets from listed files")
	fsListCmd.Flags().BoolP("random", "r", false, "Create vegeta targets from listed files with random names, works only together with -target")
}
