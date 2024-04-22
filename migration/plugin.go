package migration

import (
	"github.com/NubeIO/module-migration/utils/sshclient"
	"log"
)

func RemovePlugins(ip, sshUsername, sshPassword string) error {
	client, err := sshclient.New(ip, sshUsername, sshPassword)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if err = session.Run("rm -rf /data/rubix-os/data/plugins/*"); err != nil {
		return err
	}
	return restartROS(client)
}
