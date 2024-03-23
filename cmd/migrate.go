package cmd

import (
	"github.com/NubeIO/module-migration/migration"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "migrate",
	Short: "For applying migration",
	Long:  "For applying migration",
	Run:   migrate,
}

func migrate(_ *cobra.Command, _ []string) {
	migration.BackupAndMigrate()
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
