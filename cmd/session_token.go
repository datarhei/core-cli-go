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
		remotes, _ := cmd.Flags().GetStringSlice("remote")
		extras, _ := cmd.Flags().GetStringSlice("extra")
		ttlFlag, _ := cmd.Flags().GetString("ttl")

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
			TTL:    int64(ttl.Seconds()),
		}

		if len(extras) != 0 {
			req.Extra = map[string]interface{}{}

			for _, e := range extras {
				before, after, _ := strings.Cut(e, ":")
				req.Extra[before] = after
			}
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
	sessionTokenCmd.Flags().StringSliceP("remote", "r", []string{}, "Comma separated list of allowed referrer hosts")
	sessionTokenCmd.Flags().StringSliceP("extra", "e", []string{}, "Comma separates list of key:value extra values")
	sessionTokenCmd.Flags().StringP("ttl", "t", "24h", "Validity duration of the token")
}
