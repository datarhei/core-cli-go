package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/spf13/cobra"
)

var clusterProcessStressAddCmd = &cobra.Command{
	Use:   "add [template] [owner] [domain]?",
	Short: "Process API stress",
	Long:  "Process API stress",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		template := args[0]
		owner := args[1]
		domain := ""
		if len(args) == 4 {
			domain = args[2]
		}

		addProcess, _ := cmd.Flags().GetBool("add")
		deleteProcess, _ := cmd.Flags().GetBool("delete")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		config, err := loadTemplate(template)
		if err != nil {
			return err
		}

		config.ID = "testadd_0"
		config.Owner = owner
		config.Domain = domain

		if addProcess {
			err = client.ClusterProcessAdd(config)
			if err != nil {
				return err
			}
		}

		process, err := client.ClusterProcess(coreclient.NewProcessID(config.ID, config.Domain), []string{"state"})
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, process, true); err != nil {
			return err
		}

		if deleteProcess {
			err = client.ClusterProcessDelete(coreclient.NewProcessID(config.ID, config.Domain), true)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	clusterProcessStressCmd.AddCommand(clusterProcessStressAddCmd)

	clusterProcessStressAddCmd.Flags().BoolP("add", "a", true, "Add process")
	clusterProcessStressAddCmd.Flags().BoolP("delete", "d", true, "Delete process")
}
