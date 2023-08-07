package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

var sessionTokenCmd = &cobra.Command{
	Use:   "token [username] [match]",
	Short: "Create a session token",
	Long:  "Create a session token",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		remoteFlag, _ := cmd.Flags().GetString("remote")
		ttlFlag, _ := cmd.Flags().GetString("ttl")

		remotes := strings.Split(remoteFlag, ",")
		ttl, err := time.ParseDuration(ttlFlag)
		if err != nil {
			return err
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		req := api.SessionTokenRequest{
			Match:  args[1],
			Remote: remotes,
			Extra:  map[string]interface{}{},
			TTL:    int64(ttl.Seconds()),
		}

		token, err := client.SessionToken(args[0], []api.SessionTokenRequest{req})
		if err != nil {
			return err
		}

		if err := writeJSON(os.Stdout, token[0].Token, true); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	sessionCmd.AddCommand(sessionTokenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//processCmd.PersistentFlags().Bool("raw", false, "Display raw result from the API as JSON")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sessionTokenCmd.Flags().StringP("remote", "r", "", "Comma separated list of allowed referrer hosts")
	sessionTokenCmd.Flags().StringP("ttl", "t", "24h", "Validity duration of the token")
}
