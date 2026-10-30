package database

import (
	"fmt"
	"strconv"

	"backend/src/database"
	"backend/src/database/seeders"

	"github.com/spf13/cobra"
)

// Cmd represents the database management command group.
var Cmd = &cobra.Command{
	Use:   "database",
	Short: "Database management commands",
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.MigrateUp()
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back all migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.MigrateDown()
	},
}

var migrateForceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force set migration version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid version number: %w", err)
		}
		return database.MigrateForce(v)
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Execute database seeders",
	RunE: func(cmd *cobra.Command, args []string) error {
		return seeders.RunAll()
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateForceCmd)

	Cmd.AddCommand(migrateCmd)
	Cmd.AddCommand(seedCmd)
}
