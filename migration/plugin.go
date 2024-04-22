package migration

import (
	"github.com/NubeIO/module-migration/utils/sshclient"
	"log"
)

func RemovePlugins(ip, sshUsername, sshPassword, port string) error {
	client, err := sshclient.New(ip, sshUsername, sshPassword, port)
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

	if err = session.Run("sudo rm -rf /data/rubix-os/data/plugins/*"); err != nil {
		return err
	}
	return restartROS(client)
}
