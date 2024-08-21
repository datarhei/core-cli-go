package cmd

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List templates",
	Long:  "List templates.",
	RunE: func(cmd *cobra.Command, args []string) error {
		templates := viper.GetStringMap("templates")

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Template"})

		for key := range templates {
			t.AppendRow(table.Row{key})
		}

		t.SetAutoIndex(true)

		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Asc},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
