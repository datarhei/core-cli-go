package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// fsDeleteCmd represents the list command
var fsDeleteCmd = &cobra.Command{
	Use:   "delete [fsname] [pattern]",
	Short: "Delete files based on pattern",
	Long:  "Delete files with the given pattern from the filesystem.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		storage := args[0]
		pattern := args[1]
		execute, _ := cmd.Flags().GetBool("execute")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		files, err := client.FilesystemList(storage, pattern, "", "")
		if err != nil {
			return err
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Name", "Size", "Last Modification"})

		nfiles := len(files)
		nbytes := uint64(0)

		for _, file := range files {
			if execute {
				err := client.FilesystemDeleteFile(storage, file.Name)
				if err != nil {
					fmt.Printf("%s\n", err.Error())
				}
			}

			lastMod := time.Unix(file.LastMod, 0)
			t.AppendRow(table.Row{file.Name, formatByteCountBinary(uint64(file.Size)), lastMod.Format("2006-01-02 15:04:05")})

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
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil
	},
}

func init() {
	fsCmd.AddCommand(fsDeleteCmd)

	fsDeleteCmd.Flags().BoolP("execute", "x", false, "Actually delete the files")
}
