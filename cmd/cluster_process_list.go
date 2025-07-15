package cmd

import (
	"os"
	"strings"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"
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
		reference, _ := cmd.Flags().GetString("reference")
		idpattern, _ := cmd.Flags().GetString("idpattern")
		refpattern, _ := cmd.Flags().GetString("refpattern")
		ownerpattern, _ := cmd.Flags().GetString("ownerpattern")
		domainpattern, _ := cmd.Flags().GetString("domainpattern")
		load, _ := cmd.Flags().GetBool("load")
		sort, _ := cmd.Flags().GetString("sort")

		list, err := client.ClusterProcessList(coreclient.ProcessListOptions{
			ID:            strings.Split(ids, ","),
			Filter:        strings.Split(filter, ","),
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

		aboutv1, aboutv2, err := client.Cluster()
		if err != nil {
			return err
		}

		nodes := map[string]api.NodeResources{}

		if load {
			var about api.ClusterAbout

			if aboutv1 != nil {
				about = aboutv1.ClusterAbout
			} else {
				about = aboutv2.ClusterAbout
			}

			for _, n := range about.Nodes {
				nodes[n.ID] = n.Resources
			}
		}

		processTable(list, pmap, nodes, sort)

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessListCmd)

	clusterProcessListCmd.Flags().String("id", "", "A comma-separated list of process IDs")
	clusterProcessListCmd.Flags().String("filter", "state", "A comma-separated list of filters per process: config, state, report, metadata")
	clusterProcessListCmd.Flags().String("reference", "", "Limit list to specific reference")
	clusterProcessListCmd.Flags().String("idpattern", "", "A glob pattern for the process IDs")
	clusterProcessListCmd.Flags().String("refpattern", "", "A glob pattern for the process references")
	clusterProcessListCmd.Flags().String("ownerpattern", "", "A gob pattern for the process owners")
	clusterProcessListCmd.Flags().String("domainpattern", "", "A gob pattern for the process domains")

	clusterProcessListCmd.Flags().Bool("load", false, "Whether to show node resources")
	clusterProcessListCmd.Flags().String("sort", "", "Table sorting")
}
