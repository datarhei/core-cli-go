package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var clusterFsListCmd = &cobra.Command{
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

		list, err := client.ClusterFilesystemList(name, pattern, sort, order)
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

		totalSize := uint64(0)

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Name", "Size", "Last Modification", "Node"})

		for _, f := range list {
			lastMod := time.Unix(f.LastMod, 0)
			t.AppendRow(table.Row{f.Name, formatByteCountBinary(uint64(f.Size)), lastMod.Format("2006-01-02 15:04:05"), f.CoreID})
			totalSize += uint64(f.Size)
		}

		t.AppendFooter(table.Row{
			strconv.Itoa(len(list)),
			formatByteCountBinary(uint64(totalSize)),
		})

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 2, Align: text.AlignRight},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil

	},
}

func init() {
	clusterFsCmd.AddCommand(clusterFsListCmd)

	clusterFsListCmd.Flags().StringP("sort", "s", "none", "Sorting criteria")
	clusterFsListCmd.Flags().StringP("order", "o", "asc", "Sorting direction")
	clusterFsListCmd.Flags().StringP("target", "t", "", "Create vegeta targets from listed files")
	clusterFsListCmd.Flags().BoolP("random", "r", false, "Create vegeta targets from listed files with random names, works only together with -target")
}
