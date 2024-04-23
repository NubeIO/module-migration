package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	ip            string
	schema        string
	port          int
	externalToken string
	sshUsername   string
	sshPassword   string
	sshPort       string
)

var RootCmd = &cobra.Command{
	Use:   "migration-cli",
	Short: "Migration CLI",
	Long:  "Migration CLI",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&schema, "schema", "", "http", "schema (default http)")
	RootCmd.PersistentFlags().StringVarP(&ip, "ip", "", "0.0.0.0", "ip (default 0.0.0.0")
	RootCmd.PersistentFlags().IntVarP(&port, "port", "", 1660, "port (default 1660)")
	RootCmd.PersistentFlags().StringVarP(&externalToken, "external-token", "", "", "external token")
	RootCmd.PersistentFlags().StringVarP(&sshUsername, "ssh-username", "", "", "ssh username")
	RootCmd.PersistentFlags().StringVarP(&sshPassword, "ssh-password", "", "", "ssh password")
	RootCmd.PersistentFlags().StringVarP(&sshPort, "ssh-port", "", "22", "ssh port")
}
