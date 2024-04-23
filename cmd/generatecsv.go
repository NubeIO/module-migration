package cmd

import (
	"github.com/NubeIO/module-migration/cli"
	"github.com/NubeIO/module-migration/utils/host"
	"github.com/spf13/cobra"
)

var generateCsvCmd = &cobra.Command{
	Use:   "generate-csv",
	Short: "Generate CSV CLI",
	Long:  "Generate CSV CLI",
	Run:   generateCSV,
}

func generateCSV(_ *cobra.Command, _ []string) {
	https := schema == "https"
	cli.Setup(ip, port, &https, externalToken)

	host.GenerateHosts()
}

func init() {
	RootCmd.AddCommand(generateCsvCmd)
}
