package cmd

import (
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"
	coreclientapi "github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

// processShowCmd represents the show command
var processShowCmd = &cobra.Command{
	Use:   "show [processid]",
	Short: "Show the process with the given ID",
	Long:  "Show the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]
		asRaw, _ := cmd.Flags().GetBool("raw")
		withConfig, _ := cmd.Flags().GetBool("cfg")
		withReport, _ := cmd.Flags().GetBool("report")
		withMetadata, _ := cmd.Flags().GetBool("metadata")

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		filter := []string{"state"}
		if withConfig {
			filter = append(filter, "config")
		}
		if withReport {
			filter = append(filter, "report")
		}
		if withMetadata {
			filter = append(filter, "metadata")
		}

		p, err := client.Process(id, filter)
		if err != nil {
			return err
		}

		if asRaw {
			if err := writeJSON(os.Stdout, p, true); err != nil {
				return err
			}

			return nil
		}

		nodeid := ""
		if about, err := client.About(true); err == nil {
			nodeid = about.ID
		}

		pmap := map[string]string{}
		pmap[coreclient.NewProcessID(p.ID, p.Domain).String()] = nodeid

		processTable([]coreclientapi.Process{p}, pmap)

		processIO(p)

		return nil
	},
}

func init() {
	processCmd.AddCommand(processShowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	processShowCmd.Flags().BoolP("cfg", "c", false, "Include the process config")
	processShowCmd.Flags().BoolP("report", "r", false, "Include the process config")
	processShowCmd.Flags().BoolP("metadata", "m", false, "Include the process config")
}
