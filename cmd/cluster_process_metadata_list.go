package cmd

import (
	"fmt"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var clusterProcessMetadataListCmd = &cobra.Command{
	Use:   "list [processid]",
	Short: "List all metadata keys",
	Long:  "List all metadata keys",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		m, err := client.ClusterProcessMetadata(id, "")
		if err != nil {
			return err
		}

		metadata, ok := m.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unknown metadata format")
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Key"})

		for key := range metadata {
			t.AppendRow(table.Row{key})
		}

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil

	},
}

func init() {
	clusterProcessMetadataCmd.AddCommand(clusterProcessMetadataListCmd)
}
