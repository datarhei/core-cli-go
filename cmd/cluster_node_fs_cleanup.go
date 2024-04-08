package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var clusterNodeFilesystemCleanupCmd = &cobra.Command{
	Use:   "cleanup [storage] [pattern]?",
	Short: "Cleanup files on a filesystem",
	Long:  "Cleanup duplicate files on a filesystem in the cluster. Only the newest entry is kept.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		storage := args[0]
		pattern := ""
		execute, _ := cmd.Flags().GetBool("x")

		if len(args) > 1 {
			pattern = args[1]
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		list, err := client.ClusterFilesystemList(storage, pattern, "", "")
		if err != nil {
			return err
		}

		files := map[string][]api.FileInfo{}

		for _, file := range list {
			f, ok := files[file.Name]
			if !ok {
				f = []api.FileInfo{}
			}

			f = append(f, file)

			files[file.Name] = f
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Name", "Duplicates", "Node", "Last Modification", "Others"})

		nfiles := uint64(0)
		nduplicates := uint64(0)

		for name, list := range files {
			if len(list) == 1 {
				continue
			}

			nfiles++

			sort.Slice(list, func(i, j int) bool { return list[i].LastMod > list[j].LastMod })

			others := []string{}
			for _, file := range list[1:] {
				if execute {
					client.ClusterNodeFilesystemDeleteFile(file.CoreID, storage, file.Name)
				}
				nduplicates++
				others = append(others, file.CoreID)
			}

			lastMod := time.Unix(list[0].LastMod, 0)
			t.AppendRow(table.Row{name, len(list) - 1, list[0].CoreID, lastMod.Format("2006-01-02 15:04:05"), strings.Join(others, ",")})
		}

		if nfiles == 0 {
			fmt.Println("No duplicate files")
			return nil
		}

		t.AppendFooter(table.Row{
			strconv.FormatUint(nfiles, 10),
			strconv.FormatUint(nduplicates, 10),
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
	clusterNodeFilesystemCmd.AddCommand(clusterNodeFilesystemCleanupCmd)

	clusterNodeFilesystemCleanupCmd.Flags().BoolP("execute", "x", false, "Actually execute the cleanup")
}
