package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var clusterFilesystemDeleteCmd = &cobra.Command{
	Use:   "delete [storage] [pattern]?",
	Short: "Delete files from a filesystem",
	Long:  "Delete files from a filesystem in the cluster.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		storage := args[0]
		pattern := ""
		execute, _ := cmd.Flags().GetBool("execute")

		if len(args) > 1 {
			pattern = args[1]
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		files, err := client.ClusterFilesystemList(storage, pattern, "", "")
		if err != nil {
			return err
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Name", "Size", "Last Modification", "Node"})

		nfiles := len(files)
		nbytes := uint64(0)

		for _, file := range files {
			if execute {
				err := client.ClusterNodeFilesystemDeleteFile(file.CoreID, storage, file.Name)
				if err != nil {
					fmt.Printf("%s\n", err.Error())
				}
			}

			lastMod := time.Unix(file.LastMod, 0)
			t.AppendRow(table.Row{file.Name, formatByteCountBinary(uint64(file.Size)), lastMod.Format("2006-01-02 15:04:05"), file.CoreID})

			nbytes += uint64(file.Size)
		}

		if nfiles == 0 {
			fmt.Println("No files found")
			return nil
		}

		t.AppendFooter(table.Row{
			strconv.Itoa(nfiles),
			formatByteCountBinary(nbytes),
		})

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 2, Align: text.AlignRight},
			{Number: 3, Align: text.AlignRight},
			{Number: 4, Align: text.AlignRight},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil

	},
}

func init() {
	clusterFilesystemCmd.AddCommand(clusterFilesystemDeleteCmd)

	clusterFilesystemDeleteCmd.Flags().BoolP("execute", "x", false, "Actually delete the files")
}
