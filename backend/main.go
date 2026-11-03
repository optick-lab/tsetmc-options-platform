package main

import (
	"os"

	"backend/cmd"
	"backend/src/pkg/logger"
)

func main() {
	defer logger.Sync()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
