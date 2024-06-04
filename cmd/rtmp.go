package cmd

import (
	"fmt"
	"net"
	"net/url"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var rtmpCmd = &cobra.Command{
	Use:   "rtmp",
	Short: "RTMP related commands",
	Long:  "RTMP related commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		rtmpPort := "1935"
		rtmpToken := ""

		if version, config, err := client.Config(); err == nil {
			if version == 1 {
				cfg, ok := config.Config.(api.ConfigV1)
				if !ok {
					return fmt.Errorf("failed to convert config")
				}

				rtmpToken = cfg.RTMP.Token
				if _, port, err := net.SplitHostPort(cfg.RTMP.Address); err != nil {
					rtmpPort = port
				}
			} else if version == 2 {
				cfg, ok := config.Config.(api.ConfigV2)
				if !ok {
					return fmt.Errorf("failed to convert config")
				}

				rtmpToken = cfg.RTMP.Token
				if _, port, err := net.SplitHostPort(cfg.RTMP.Address); err != nil {
					rtmpPort = port
				}
			} else if version == 3 {
				cfg, ok := config.Config.(api.ConfigV3)
				if !ok {
					return fmt.Errorf("failed to convert config")
				}

				rtmpToken = cfg.RTMP.Token
				if _, port, err := net.SplitHostPort(cfg.RTMP.Address); err != nil {
					rtmpPort = port
				}
			}
		} else {
			return err
		}

		channels, err := client.RTMPChannels()
		if err != nil {
			return err
		}

		t := table.NewWriter()

		t.AppendHeader(table.Row{"Channel", "URL"})

		for _, channel := range channels {
			u, err := url.Parse(client.Address())
			if err == nil {
				u.Scheme = "rtmp"
				u.Host = u.Hostname() + ":" + rtmpPort
				u.Path = channel.Name

				if len(rtmpToken) != 0 {
					u = u.JoinPath(rtmpToken)
				}
			}

			t.AppendRow(table.Row{
				channel.Name,
				u.String(),
			})
		}

		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Asc},
		})

		t.SetStyle(table.StyleLight)

		fmt.Println(t.Render())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rtmpCmd)
}
