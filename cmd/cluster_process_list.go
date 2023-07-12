package cmd

import (
	"os"
	"strings"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/spf13/cobra"
)

var clusterProcessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all processes",
	Long:  "List all processes in the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		asRaw, _ := cmd.Flags().GetBool("raw")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		ids, _ := cmd.Flags().GetString("ids")
		filter, _ := cmd.Flags().GetString("filter")
		domain, _ := cmd.Flags().GetString("domain")
		reference, _ := cmd.Flags().GetString("reference")
		idpattern, _ := cmd.Flags().GetString("idpattern")
		refpattern, _ := cmd.Flags().GetString("refpattern")
		ownerpattern, _ := cmd.Flags().GetString("ownerpattern")
		domainpattern, _ := cmd.Flags().GetString("domainpattern")

		list, err := client.ClusterProcessList(coreclient.ProcessListOptions{
			ID:            strings.Split(ids, ","),
			Filter:        strings.Split(filter, ","),
			Domain:        domain,
			Reference:     reference,
			IDPattern:     idpattern,
			RefPattern:    refpattern,
			OwnerPattern:  ownerpattern,
			DomainPattern: domainpattern,
		})
		if err != nil {
			return err
		}

		if asRaw {
			if err := writeJSON(os.Stdout, list, true); err != nil {
				return err
			}

			return nil
		}

		pmap, err := client.ClusterDBProcessMap()
		if err != nil {
			return err
		}

		processTable(list, pmap)

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessListCmd)

	clusterProcessListCmd.Flags().String("id", "", "A comma-separated list of process IDs")
	clusterProcessListCmd.Flags().String("filter", "state", "A comma-separated list of filters per process: config, state, report, metadata")
	clusterProcessListCmd.Flags().String("domain", "", "The domain to act upon")
	clusterProcessListCmd.Flags().String("reference", "", "Limit list to specific reference")
	clusterProcessListCmd.Flags().String("idpattern", "", "A glob pattern for the process IDs")
	clusterProcessListCmd.Flags().String("refpattern", "", "A glob pattern for the process references")
	clusterProcessListCmd.Flags().String("ownerpattern", "", "A gob pattern for the process owners")
	clusterProcessListCmd.Flags().String("domainpattern", "", "A gob pattern for the process domains")
}
