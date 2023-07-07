package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type coreEdit struct {
	Address      string `json:"address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

var coreEditCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit a core",
	Long:  `Edit a core.`,
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := viper.GetString("cores.selected")
		if len(globalFlagCore) != 0 {
			name = globalFlagCore
		}

		if len(args) > 0 {
			name = args[0]
		}

		list := viper.GetStringMapString("cores.list")

		coreURL, ok := list[name]
		if !ok {
			return fmt.Errorf("core with name '%s' not found", name)
		}

		u, err := url.Parse(coreURL)
		if err != nil {
			return fmt.Errorf("invalid data for core '%s': %w", name, err)
		}

		password, _ := u.User.Password()
		query := u.Query()

		core := coreEdit{
			Address:      u.Scheme + "://" + u.Host + u.Path,
			Username:     u.User.Username(),
			Password:     password,
			AccessToken:  u.Query().Get("accessToken"),
			RefreshToken: u.Query().Get("refreshToken"),
		}

		data, err := json.MarshalIndent(core, "", "   ")
		if err != nil {
			return err
		}

		editedData, modified, err := editData(data)
		if err != nil {
			return err
		}

		if !modified {
			// They are the same, nothing has been changed. No need to store the metadata
			fmt.Printf("No changes. Core config will not be updated.")
			return nil
		}

		editedCore := coreEdit{}

		if err := json.Unmarshal(editedData, &editedCore); err != nil {
			return err
		}

		u, err = url.Parse(editedCore.Address)
		if err != nil {
			return err
		}

		if len(editedCore.Username) != 0 {
			if len(password) == 0 {
				u.User = url.User(editedCore.Username)
			} else {
				u.User = url.UserPassword(editedCore.Username, editedCore.Password)
			}
		}

		u.User = url.UserPassword(editedCore.Username, editedCore.Password)

		if len(editedCore.AccessToken) == 0 {
			query.Del("accessToken")
		} else {
			query.Set("accessToken", editedCore.AccessToken)
		}

		if len(editedCore.RefreshToken) == 0 {
			query.Del("refreshToken")
		} else {
			query.Set("refreshToken", editedCore.RefreshToken)
		}

		u.RawQuery = query.Encode()

		list[name] = u.String()

		viper.Set("cores.list", list)
		viper.WriteConfig()

		return nil
	},
}

func init() {
	coreCmd.AddCommand(coreEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// selectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// selectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
