package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var clusterProcessProbeCmd = &cobra.Command{
	Use:   "probe [processid]?",
	Short: "Probe the process with the given ID",
	Long:  "Probe the process with the given ID",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		var probe api.Probe

		if len(args) == 1 {
			pid := args[0]

			probe, err = clusterProbeFromID(client, pid)
		} else {
			fromFile, _ := cmd.Flags().GetString("from-file")
			if len(fromFile) == 0 {
				return fmt.Errorf("no process configuration file provided")
			}

			coreid, _ := cmd.Flags().GetString("coreid")

			reader := os.Stdin

			if fromFile != "-" {
				if file, err := os.Open(fromFile); err != nil {
					return err
				} else {
					reader = file
				}
			}

			probe, err = clusterProbeFromConfig(client, reader, coreid)
		}

		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, probe, true); err != nil {
			return err
		}

		return nil
	},
}

func clusterProbeFromID(client coreclient.RestClient, pid string) (api.Probe, error) {
	id := coreclient.ParseProcessID(pid)

	return client.ClusterProcessProbe(id)
}

func clusterProbeFromConfig(client coreclient.RestClient, r io.Reader, coreid string) (api.Probe, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return api.Probe{}, err
	}

	config := api.ProcessConfig{}

	if err := json.Unmarshal(data, &config); err != nil {
		return api.Probe{}, err
	}

	return client.ClusterProcessProbeConfig(config, coreid)
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessProbeCmd)

	clusterProcessProbeCmd.Flags().String("from-file", "-", "Load process config from file or stdin")
	clusterProcessProbeCmd.Flags().String("coreid", "", "Execute the probe preferably on this core")
}
