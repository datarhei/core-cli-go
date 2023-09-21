package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildkite/shellwords"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start a REPL",
	Long:  "Start a REPL",
	RunE: func(cmd *cobra.Command, args []string) error {
		line := liner.NewLiner()
		defer line.Close()

		line.SetCtrlCAborts(false)

		home, err := os.UserHomeDir()
		if err != nil {
			home = "./"
		}
		historyFilepath := filepath.Join(home, ".corecli.history")

		if f, err := os.Open(historyFilepath); err == nil {
			line.ReadHistory(f)
			f.Close()
		}

		for {
			selected := viper.GetString("cores.selected")
			if command, err := line.Prompt(fmt.Sprintf("%s> ", selected)); err == nil {
				if strings.ToLower(command) == "exit" {
					break
				}

				if len(command) == 0 {
					continue
				}

				args, err := shellwords.Split(command)
				if err != nil {
					fmt.Printf("Error reading line: %s\n", err.Error())
					continue
				}

				line.AppendHistory(command)

				rootCmd.SetArgs(args)
				rootCmd.Execute()
			} else {
				fmt.Printf("Error reading line: %s\n", err.Error())
			}
		}

		if f, err := os.Create(historyFilepath); err == nil {
			line.WriteHistory(f)
			f.Close()
		} else {
			fmt.Printf("Failed storing history: %s\n", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}
