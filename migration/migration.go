package migration

import (
	"fmt"
	"github.com/NubeIO/flow-eng/helpers/boolean"
	"github.com/NubeIO/module-migration/cli"
	"github.com/NubeIO/module-migration/utils/host"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
)

func Migrate(sshUsername, sshPassword, sshPort string) {
	hosts, err := host.GetHosts()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	if len(hosts) <= 1 {
		fmt.Printf("Host not found.")
	}

	for _, hos := range hosts {
		if hos.LocationUUID == "Location UUID" {
			continue
		}
		port, err := strconv.Atoi(hos.Port)
		if err != nil {
			log.Fatalf("port %s is not convertible to integer", hos.Port)
		}
		cli.Setup(hos.IP, port, boolean.NewFalse(), "") // we don't need external token to check ping
		_, pingable, _ := cli.CLIShort.Ping()
		if pingable {
			if hos.PluginDeletionState != "true" {
				log.Printf("Remove plugins started for host: %s", hos.HostName)
				err = RemovePlugins(hos.IP, sshUsername, sshPassword, sshPort)
				if err == nil {
					hos.PluginDeletionState = "true"
					hos.PluginDeletionStatus = ""
				} else {
					hos.PluginDeletionStatus = err.Error()
				}
				log.Printf("Remove plugins finished for host: %s", hos.HostName)
			}

			if hos.RosMigrationState != "true" {
				log.Printf("Ros migration started for host: %s", hos.HostName)
				err = BackupAndMigrateROS(hos.IP, sshUsername, sshPassword, sshPort)
				if err == nil {
					hos.RosMigrationState = "true"
					hos.RosMigrationStatus = ""
				} else {
					hos.RosMigrationStatus = err.Error()
				}
				log.Printf("Ros migration finished for host: %s", hos.HostName)
			}

			if hos.WiresMigrationState != "true" {
				log.Printf("Wires migration started for host: %s", hos.HostName)
				err = MigrateWires(hos.IP, sshUsername, sshPassword, sshPort)
				if err == nil {
					hos.WiresMigrationState = "true"
					hos.WiresMigrationStatus = ""
				} else {
					hos.WiresMigrationStatus = err.Error()
				}
				log.Printf("Wires migration finished for host: %s", hos.HostName)
			}
		} else {
			if hos.PluginDeletionState != "true" {
				hos.PluginDeletionStatus = "device is offline"
			}

			if hos.RosMigrationState != "true" {
				hos.RosMigrationStatus = "device is offline"
			}

			if hos.WiresMigrationState != "true" {
				hos.WiresMigrationStatus = "device is offline"
			}
		}
	}

	if err = host.UpdateHosts(hosts); err != nil {
		fmt.Printf(err.Error())
	}
}

func backup(client *ssh.Client, destination, destinationDir string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo mkdir -p %s && sudo cp %s %s", destinationDir, rosDbFile, destination)
	return session.Run(cmd)
}

func giveFilePermission(client *ssh.Client, file string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo chmod 777 %s", file)
	return session.Run(cmd)
}
