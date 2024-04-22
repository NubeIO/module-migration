package migration

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-eng/nodes"
	"github.com/NubeIO/module-migration/utils/sshclient"
	"github.com/NubeIO/module-migration/utils/wiresnew"
	"github.com/NubeIO/module-migration/utils/wiresold"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var wiresDbFile = "/data/rubix-edge-wires/data/data.db"
var wiresJsonFile = "/data/rubix-edge-wires/data/data.json"
var wiresDownloadDbFile = "./wires-data.db"
var wiresDownloadJsonFile = "./data.json"

func MigrateWires(ip, sshUsername, sshPassword, sshPort string) error {
	client, err := sshclient.New(ip, sshUsername, sshPassword, sshPort)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer client.Close()

	currentDateTime := time.Now().UTC().Format("20060102150405")
	destinationDir := fmt.Sprintf("/data/backup/migration/rubix-edge-wires/%s", currentDateTime)
	destination := filepath.Join(destinationDir, "data.db")
	if err = backup(client, destination, destinationDir); err != nil {
		return fmt.Errorf("error on doing ROS backup: %s", err.Error())
	}

	_ = giveFilePermission(client, wiresDbFile)
	_ = giveFilePermission(client, wiresJsonFile)

	log.Printf("wires migration started")
	log.Printf("Started download")
	if err = downloadWiresDb(ip, sshUsername, sshPassword, sshPort); err != nil {
		return fmt.Errorf("failed to download wires: %s", err)
	}
	log.Printf("Finished download")

	nodeList, hostUUID, err := wiresold.Get(wiresDownloadDbFile)
	if err != nil {
		return fmt.Errorf("failed to read wires on old format: %s", err)
	}
	log.Printf("read wires into old wires")

	var encodedNodes nodes.NodesList
	if err = json.Unmarshal(nodeList, &encodedNodes); err != nil {
		return fmt.Errorf("failed to unmarshal wires: %s", err)
	}
	log.Printf("wires unmarshalled into nodes")

	_ = os.Remove(wiresDownloadJsonFile) // remove it otherwise it gets appended
	if err = wiresnew.Migrate(wiresDownloadJsonFile, &wiresnew.FlowDownload{
		HostUUID:     hostUUID,
		EncodedNodes: &encodedNodes,
	}); err != nil {
		return fmt.Errorf("failed to migrate wires: %s", err)
	}
	log.Printf("wires migrated into new wires")

	// Stop it first, otherwise it gets runtime values
	if err = stopWires(client); err != nil {
		return fmt.Errorf("failed to stop wires: %s", err)
	}

	log.Printf("Started upload")
	if err = uploadWiresDb(ip, sshUsername, sshPassword, sshPort); err != nil {
		return fmt.Errorf("failed to upload wires: %s", err)
	}
	log.Printf("Finished upload")

	return restartWires(client)
}

func downloadWiresDb(ip, sshUsername, sshPassword, sshPort string) error {
	cmd := fmt.Sprintf("sshpass -p '%s' scp -o StrictHostKeyChecking=no -P %s %s@%s:%s %s",
		sshPassword, sshPort, sshUsername, ip, wiresDbFile, wiresDownloadDbFile)
	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func uploadWiresDb(ip, sshUsername, sshPassword, sshPort string) error {
	cmd := fmt.Sprintf("sshpass -p '%s' scp -o StrictHostKeyChecking=no -P %s %s %s@%s:%s",
		sshPassword, sshPort, wiresDownloadJsonFile, sshUsername, ip, wiresJsonFile)
	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func restartWires(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run("sudo systemctl restart nubeio-rubix-edge-wires.service")
}

func stopWires(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run("sudo systemctl stop nubeio-rubix-edge-wires.service")
}
