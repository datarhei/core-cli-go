package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var fsRestoreCmd = &cobra.Command{
	Use:   "restore [fsname] [sourcedir]",
	Short: "Restore files",
	Long:  "Restore files from filesystem, the target filesystem will not be wiped.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		sourcedir := args[1]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		sourcedir, err = filepath.Abs(sourcedir)
		if err != nil {
			return err
		}

		filelist := map[string]uint64{}

		err = filepath.Walk(sourcedir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			mode := info.Mode()
			if !mode.IsRegular() {
				return nil
			}

			if mode&os.ModeSymlink != 0 {
				return nil
			}

			filelist[path] = uint64(info.Size())

			return nil
		})
		if err != nil {
			return err
		}

		if len(filelist) == 0 {
			return nil
		}

		totalSize, restoreSize := uint64(0), uint64(0)

		for _, size := range filelist {
			totalSize += size
		}

		fmt.Printf("Restore size: %d bytes\n", totalSize)

		start := time.Now()

		for path, size := range filelist {
			elapsed := time.Since(start).Seconds()
			fmt.Printf("%3d%% done (%12d/%12d bytes, %.0f bytes/s)\r", uint64(float64(restoreSize)/float64(totalSize)*100), restoreSize, totalSize, float64(restoreSize)/elapsed)

			file, err := os.Open(path)
			if err != nil {
				return nil
			}

			err = client.FilesystemAddFile(name, strings.TrimPrefix(path, sourcedir), file)

			file.Close()

			if err != nil {
				return err
			}

			restoreSize += size
		}

		elapsed := time.Since(start).Seconds()
		fmt.Printf("100%% done (%12d/%12d bytes, %.0f bytes/s)\n", restoreSize, totalSize, float64(restoreSize)/elapsed)

		return nil

	},
}

func init() {
	fsCmd.AddCommand(fsRestoreCmd)
}
