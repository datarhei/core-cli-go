package cmd

import (
	"github.com/spf13/cobra"
)

// skillsReloadCmd represents the list command
var skillsReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload skills",
	Long:  "Reload FFmpeg skills",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		return client.SkillsReload()
	},
}

func init() {
	skillsCmd.AddCommand(skillsReloadCmd)
}
