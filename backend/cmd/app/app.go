package app

import (
	"backend/src/bootstrap"

	"github.com/spf13/cobra"
)

// Cmd represents the app command group.
var Cmd = &cobra.Command{
	Use:   "app",
	Short: "Application lifecycle commands",
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Run database migrations and start the HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return bootstrap.Bootstrap()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server without running migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return bootstrap.Serve()
	},
}

func init() {
	Cmd.AddCommand(bootstrapCmd)
	Cmd.AddCommand(serveCmd)
}
