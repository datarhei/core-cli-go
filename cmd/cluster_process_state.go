package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/itchyny/gojq"
	"github.com/spf13/cobra"
)

var clusterProcessStateCmd = &cobra.Command{
	Use:   "state [processid]",
	Short: "Show the state of the process with the given ID",
	Long:  "Show the state of the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jq, _ := cmd.Flags().GetString("jq")
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		process, err := client.ClusterProcess(id, []string{"state"})
		if err != nil {
			return err
		}

		if len(jq) != 0 {
			query, err := gojq.Parse(jq)
			if err != nil {
				return err
			}

			data, err := json.Marshal(process.State)
			if err != nil {
				return err
			}

			var input any

			err = json.Unmarshal(data, &input)
			if err != nil {
				return err
			}

			iter := query.Run(input)
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					return err
				}
				fmt.Printf("%#v\n", v)
			}
		} else {

			if err := writeJSON(os.Stdout, process.State, true); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessStateCmd)

	clusterProcessStateCmd.Flags().String("jq", "", "Transform result with jq")
}
