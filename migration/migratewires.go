package migration

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-eng/nodes"
	"github.com/NubeIO/module-migration/utils/wiresnew"
	"github.com/NubeIO/module-migration/utils/wiresold"
	"log"
	"os/exec"
)

var wiresDbFile = "/data/rubix-edge-wires/data/data.db"
var wiresJsonFile = "/data/rubix-edge-wires/data/data.json"
var wiresDownloadDbFile = "/data/download/rubix-edge-wires/data/data.db"
var wiresDownloadJsonFile = "/data/download/rubix-edge-wires/data/data.json"

func MigrateWires(ip, sshUsername, sshPassword string) error {
	if err := downloadWiresDb(ip, sshUsername, sshPassword); err != nil {
		log.Printf(err.Error())
		return err
	}

	nodeList, hostUUID, err := wiresold.Get(wiresDownloadDbFile)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	var encodedNodes nodes.NodesList
	if err = json.Unmarshal(nodeList, &encodedNodes); err != nil {
		log.Printf(err.Error())
		return err
	}

	if err = wiresnew.Migrate(wiresDownloadJsonFile, &wiresnew.FlowDownload{
		HostUUID:     hostUUID,
		EncodedNodes: &encodedNodes,
	}); err != nil {
		log.Printf(err.Error())
		return err
	}

	if err = uploadWiresDb(ip, sshUsername, sshPassword); err != nil {
		log.Printf(err.Error())
		return err
	}

	return restartWires(ip, sshUsername, sshPassword)
}

func downloadWiresDb(ip, sshUsername, sshPassword string) error {
	cmd := fmt.Sprintf("sshpass -p '%s' scp -o StrictHostKeyChecking=no %s@%s:%s %s",
		sshPassword, sshUsername, ip, wiresDbFile, wiresDownloadDbFile)
	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func uploadWiresDb(ip, sshUsername, sshPassword string) error {
	cmd := fmt.Sprintf("sshpass -p '%s' scp -o StrictHostKeyChecking=no %s %s@%s:%s",
		sshPassword, wiresDownloadJsonFile, sshUsername, ip, wiresJsonFile)
	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func restartWires(ip, sshUsername, sshPassword string) error {
	cmd := fmt.Sprintf("sshpass -p '%s' ssh -o StrictHostKeyChecking=no %s@%s sudo systemctl restart nubeio-rubix-edge-wires.service",
		sshPassword, sshUsername, ip)
	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}
