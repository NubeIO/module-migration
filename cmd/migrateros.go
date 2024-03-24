package cmd

import (
	"github.com/NubeIO/module-migration/migration"
	"github.com/spf13/cobra"
)

var migrateROSCmd = &cobra.Command{
	Use:   "migrate-ros",
	Short: "For applying migration",
	Long:  "For applying migration",
	Run:   migrateROS,
}

func migrateROS(_ *cobra.Command, _ []string) {
	migration.BackupAndMigrateROS()
}

func init() {
	RootCmd.AddCommand(migrateROSCmd)
}
