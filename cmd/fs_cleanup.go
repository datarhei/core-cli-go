package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var fsCleanupCmd = &cobra.Command{
	Use:   "cleanup [name] [core] [pattern]?",
	Short: "cleanup files",
	Long:  "cleanup files on filesystem",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		core := args[1]
		pattern := args[2]

		client1, err := connectSelectedCore()
		if err != nil {
			return err
		}

		client2, err := connectCore(core)
		if err != nil {
			return err
		}

		list1, err := client1.FilesystemList(name, pattern, "", "")
		if err != nil {
			return err
		}

		list2, err := client2.FilesystemList(name, pattern, "", "")
		if err != nil {
			return err
		}

		for _, f1 := range list1 {
			for _, f2 := range list2 {
				if f1.Name != f2.Name {
					continue
				}

				if f1.LastMod > f2.LastMod {
					fmt.Printf("delete %s on %s\n", f2.Name, client2.Address())
					if err := client2.FilesystemDeleteFile(name, f2.Name); err != nil {
						return err
					}
				} else {
					fmt.Printf("delete %s on %s\n", f1.Name, client1.Address())
					if err := client1.FilesystemDeleteFile(name, f1.Name); err != nil {
						return err
					}
				}
			}
		}

		return nil

	},
}

func init() {
	fsCmd.AddCommand(fsCleanupCmd)
}
