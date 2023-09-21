package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var fsBackupCmd = &cobra.Command{
	Use:   "backup [fsname] [pattern|pattern|pattern|...] [targetdir]",
	Short: "Backup files",
	Long:  "Backup files on filesystem, the targetdir will be wiped.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		patterns := strings.Split(args[1], "|")
		targetdir := args[2]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		targetdir, err = filepath.Abs(targetdir)
		if err != nil {
			return err
		}

		if err := os.RemoveAll(targetdir); err != nil {
			return err
		}

		if err := os.MkdirAll(targetdir, 0755); err != nil {
			return err
		}

		filelist := map[string]uint64{}

		for _, pattern := range patterns {
			files, err := client.FilesystemList(name, pattern, "", "")
			if err != nil {
				return err
			}

			for _, file := range files {
				filelist[file.Name] = uint64(file.Size)
			}
		}

		if len(filelist) == 0 {
			return nil
		}

		totalSize, backupSize := uint64(0), uint64(0)

		for _, size := range filelist {
			totalSize += size
		}

		fmt.Printf("Backup size: %d bytes\n", totalSize)

		start := time.Now()

		for path, size := range filelist {
			elapsed := time.Since(start).Seconds()
			fmt.Printf("%3d%% done (%12d/%12d bytes, %.0f bytes/s)\r", uint64(float64(backupSize)/float64(totalSize)*100), backupSize, totalSize, float64(backupSize)/elapsed)

			file, err := client.FilesystemGetFile(name, path)
			if err != nil {
				return err
			}

			path = filepath.Join(targetdir, path)

			dir := filepath.Dir(path)
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				file.Close()
				return err
			}

			f, err := os.Create(path)
			if err != nil {
				file.Close()
				return err
			}

			_, err = f.ReadFrom(file)
			if err != nil {
				f.Close()
				file.Close()

				return err
			}

			f.Close()
			file.Close()

			backupSize += size
		}

		elapsed := time.Since(start).Seconds()
		fmt.Printf("100%% done (%12d/%12d bytes, %.0f bytes/s)\n", backupSize, totalSize, float64(backupSize)/elapsed)

		return nil

	},
}

func init() {
	fsCmd.AddCommand(fsBackupCmd)
}
