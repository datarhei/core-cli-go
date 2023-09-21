package cmd

import (
	"fmt"
	"strconv"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/spf13/cobra"
)

// processAddCmd represents the add command
var processTestDeleteCmd = &cobra.Command{
	Use:   "delete [number of processes]",
	Short: "Delete processes",
	Long:  "Delete processes.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			fmt.Printf("%4d / %4d done\r", i+1, n)

			id := "processTest-" + strconv.Itoa(i)

			if err := client.ClusterProcessDelete(coreclient.NewProcessID(id, "")); err != nil {
				return err
			}
		}

		fmt.Printf("%4d / %4d done\r", n, n)

		return nil
	},
}

func init() {
	processTestCmd.AddCommand(processTestDeleteCmd)

	//processAddCmd.Flags().String("from-file", "-", "Load process config from file or stdin")
}
