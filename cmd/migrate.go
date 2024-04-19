package cmd

import (
	"github.com/NubeIO/module-migration/migration"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "For applying migration",
	Long:  "For applying migration",
	Run:   migrate,
}

func migrate(_ *cobra.Command, _ []string) {
	migration.Migrate(sshUsername, sshPassword)
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
