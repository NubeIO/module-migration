package migration

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-eng/nodes"
	"github.com/NubeIO/module-migration/utils/wiresnew"
	"github.com/NubeIO/module-migration/utils/wiresold"
	"log"
	"os"
	"os/exec"
)

var wiresDbFile = "/data/rubix-edge-wires/data/data.db"
var wiresJsonFile = "/data/rubix-edge-wires/data/data.json"
var wiresDownloadDbFile = "./wires-data.db"
var wiresDownloadJsonFile = "./data.json"

func MigrateWires(ip, sshUsername, sshPassword string) error {
	log.Printf("wires migration started")
	if err := downloadWiresDb(ip, sshUsername, sshPassword); err != nil {
		return fmt.Errorf("failed to download wires: %s", err)
	}
	log.Printf("wires downloaded")

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
	if err = stopWires(ip, sshUsername, sshPassword); err != nil {
		return fmt.Errorf("failed to stop wires: %s", err)
	}

	if err = uploadWiresDb(ip, sshUsername, sshPassword); err != nil {
		return fmt.Errorf("failed to upload wires: %s", err)
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

func stopWires(ip, sshUsername, sshPassword string) error {
	cmd := fmt.Sprintf("sshpass -p '%s' ssh -o StrictHostKeyChecking=no %s@%s sudo systemctl stop nubeio-rubix-edge-wires.service",
		sshPassword, sshUsername, ip)
	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}
