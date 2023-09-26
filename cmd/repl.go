package cmd

import (
	"fmt"
	"io"
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

		stdout := os.Stdout
		stdin := os.Stdin
		var redirectReader, redirectWriter *os.File
		redirect := ""
		redirectArgs := []string{}

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

				for i, arg := range args {
					if arg == ">" {
						redirect = "stdout"
						args, redirectArgs = args[:i], args[i:]
						break
					} else if arg == "<" {
						redirect = "stdin"
						args, redirectArgs = args[:i], args[i:]
						break
					}
				}

				if len(redirect) != 0 && len(redirectArgs) != 2 {
					fmt.Printf("Redirect required exactly one argument\n")
					continue
				}

				if redirect == "stdout" {
					file, err := os.Create(redirectArgs[1])
					if err != nil {
						fmt.Printf("Error creating output file: %s\n", err.Error())
						continue
					}

					redirectReader, redirectWriter, err = os.Pipe()
					if err != nil {
						fmt.Printf("Error redirecting output: %s\n", err.Error())
						file.Close()
					}
					os.Stdout = redirectWriter

					go func(file *os.File) {
						io.Copy(file, redirectReader)
						file.Close()
						redirectReader.Close()
					}(file)
				} else if redirect == "stdin" {
					file, err := os.Open(redirectArgs[1])
					if err != nil {
						fmt.Printf("Error opening input file: %s\n", err.Error())
						continue
					}

					redirectReader, redirectWriter, err = os.Pipe()
					if err != nil {
						fmt.Printf("Error redirecting input: %s\n", err.Error())
						file.Close()
					}
					os.Stdin = redirectReader

					go func(file *os.File) {
						io.Copy(redirectWriter, file)
						file.Close()
						redirectWriter.Close()
					}(file)
				}

				line.AppendHistory(command)

				rootCmd.SetArgs(args)
				rootCmd.Execute()

				if redirect == "stdout" {
					redirectWriter.Close()
					os.Stdout = stdout
				} else if redirect == "stdin" {
					redirectReader.Close()
					os.Stdin = stdin
				}
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
