package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var clusterNodeFilesystemListCmd = &cobra.Command{
	Use:   "list [nodeid] [fsname] [pattern]? (-s|--sort) [none|name|size|lastmod] (-o|--order) [asc|desc] (-t|--target) [url with %%s]",
	Short: "List files",
	Long:  "List files on filesystem",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		storage := args[1]
		pattern := ""
		if len(args) == 3 {
			pattern = args[2]
		}

		sort, _ := cmd.Flags().GetString("sort")
		order, _ := cmd.Flags().GetString("order")
		target, _ := cmd.Flags().GetString("target")
		random, _ := cmd.Flags().GetBool("random")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ClusterNodeFilesystemList(id, storage, pattern, sort, order)
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
	clusterNodeFilesystemCmd.AddCommand(clusterNodeFilesystemListCmd)

	clusterNodeFilesystemListCmd.Flags().StringP("sort", "s", "none", "Sorting criteria")
	clusterNodeFilesystemListCmd.Flags().StringP("order", "o", "asc", "Sorting direction")
	clusterNodeFilesystemListCmd.Flags().StringP("target", "t", "", "Create vegeta targets from listed files")
	clusterNodeFilesystemListCmd.Flags().BoolP("random", "r", false, "Create vegeta targets from listed files with random names, works only together with -target")
}
