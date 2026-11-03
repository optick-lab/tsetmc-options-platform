package cmd

import (
	"backend/cmd/app"
	"backend/cmd/database"
	"backend/src/config"
	"backend/src/pkg/logger"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backend",
	Short: "TSETMC Options Platform Backend Service",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfigAndLogger)

	rootCmd.AddCommand(app.Cmd)
	rootCmd.AddCommand(database.Cmd)
}

func initConfigAndLogger() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)
}
