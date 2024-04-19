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
	if len(hosts) == 0 {
		fmt.Printf("Host not found.")
	}

	for _, hos := range hosts {
		if !hos.PluginDeletion {
			log.Printf("Remove plugins started for host: %s", hos.HostName)
			err = RemovePlugins(hos.VirtualIP, sshUsername, sshPassword)
			if err == nil {
				hos.PluginDeletion = true
			}
			log.Printf("Remove plugins finished for host: %s", hos.HostName)
		}

		if !hos.RosMigration {
			log.Printf("Ros migration started for host: %s", hos.HostName)
			err = BackupAndMigrateROS(hos.VirtualIP, sshUsername, sshPassword)
			if err == nil {
				hos.RosMigration = true
			}
			log.Printf("Ros migration finished for host: %s", hos.HostName)
		}

		if !hos.WiresMigration {
			log.Printf("Wires migration started for host: %s", hos.HostName)
			err = MigrateWires(hos.VirtualIP, sshUsername, sshPassword)
			if err == nil {
				hos.WiresMigration = true
			}
			log.Printf("Wires migration finished for host: %s", hos.HostName)
		}
	}

	if err = host.UpdateHosts(hosts); err != nil {
		fmt.Printf(err.Error())
	}
}
