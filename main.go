package main

import (
	"os"

	"github.com/datarhei/core-cli-go/cmd"
)

func main() {
	// If there are no command line args, start in REPL mode
	if len(os.Args) == 1 {
		cmd.SetArgs([]string{"repl"})
	}

	cmd.Execute()
}
