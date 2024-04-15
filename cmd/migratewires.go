package cmd

import (
	"github.com/NubeIO/module-migration/migration"
	"github.com/spf13/cobra"
)

var migrateWiresCmd = &cobra.Command{
	Use:   "migrate-wires",
	Short: "For applying migration",
	Long:  "For applying migration",
	Run:   migrateWires,
}

func migrateWires(_ *cobra.Command, _ []string) {
	migration.BackupWires()
	migration.MigrateWires()
}

func init() {
	RootCmd.AddCommand(migrateWiresCmd)
}
