package migration

import (
	"fmt"
	"github.com/NubeIO/module-migration/utils/host"
	"log"
)

func Migrate(sshUsername, sshPassword string) {
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
		if hos.PluginDeletionState != "true" {
			log.Printf("Remove plugins started for host: %s", hos.HostName)
			err = RemovePlugins(hos.IP, sshUsername, sshPassword)
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
			err = BackupAndMigrateROS(hos.IP, sshUsername, sshPassword)
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
			err = MigrateWires(hos.IP, sshUsername, sshPassword)
			if err == nil {
				hos.WiresMigrationState = "true"
				hos.WiresMigrationStatus = ""
			} else {
				hos.WiresMigrationStatus = err.Error()
			}
			log.Printf("Wires migration finished for host: %s", hos.HostName)
		}
	}

	if err = host.UpdateHosts(hosts); err != nil {
		fmt.Printf(err.Error())
	}
}
