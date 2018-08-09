package main

import (
	"os"

	"github.com/Manjukb/enduserbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
