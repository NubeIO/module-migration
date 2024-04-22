package migration

import (
	"errors"
	"fmt"
	"github.com/NubeIO/module-migration/utils/sshclient"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/ssh"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var rosDbFile = "/data/rubix-os/data/data.db"
var rosDownloadDbFile = "./ros-data.db"

func BackupAndMigrateROS(ip, sshUsername, sshPassword, sshPort string) error {
	client, err := sshclient.New(ip, sshUsername, sshPassword, sshPort)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer client.Close()

	currentDateTime := time.Now().UTC().Format("20060102150405")
	destinationDir := fmt.Sprintf("/data/backup/migration/rubix-os/%s", currentDateTime)
	destination := filepath.Join(destinationDir, "data.db")
	if err = backupROS(client, destination, destinationDir); err != nil {
		return fmt.Errorf("error on doing ROS backup: %s", err.Error())
	}

	_ = giveFilePermission(client, rosDbFile)

	log.Printf("Started download")
	if err := downloadRosDb(ip, sshUsername, sshPassword, sshPort); err != nil {
		return fmt.Errorf("error on downloading ROS DB: %s", err.Error())
	}
	log.Printf("Finished download")

	log.Printf("Started migrating ROS data")
	if err = migrateROSData(); err != nil {
		_ = restartROS(client)
		return fmt.Errorf("error on doing ROS migration: %s", err.Error())
	}
	log.Printf("Finished migrating ROS data")

	log.Printf("Started upload")
	if err := uploadRosDb(ip, sshUsername, sshPassword, sshPort); err != nil {
		return fmt.Errorf("error on uploading ROS DB: %s", err.Error())
	}
	log.Printf("Finished upload")
	return restartROS(client)
}

func backupROS(client *ssh.Client, destination, destinationDir string) error {
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

func migrateROSData() error {
	cmd := `'SELECT COUNT(*) FROM networks WHERE plugin_name = "lora"'`
	out, err := runSqliteCommand(cmd)
	if err != nil {
		return err
	}
	count, _ := strconv.Atoi(strings.TrimRight(string(out), "\n"))
	if count > 0 {
		log.Printf("loraraw exist")
		cmd = `'SELECT COUNT(*) FROM plugins WHERE name = "module-core-loraraw"'`
		out, err = runSqliteCommand(cmd)
		if err != nil {
			return err
		}
		count, _ = strconv.Atoi(strings.TrimRight(string(out), "\n"))
		if count != 1 {
			return errors.New("ERROR: install module-core-loraraw at first")
		}
	} else {
		log.Printf("loraraw doesn't exist")
	}

	cmd = `'SELECT COUNT(*) FROM networks WHERE plugin_name = "lorawan"'`
	out, err = runSqliteCommand(cmd)
	if err != nil {
		return err
	}
	count, _ = strconv.Atoi(strings.TrimRight(string(out), "\n"))
	if count > 0 {
		log.Printf("lorawan exist")
		cmd = `'SELECT COUNT(*) FROM plugins WHERE name = "module-core-lorawan"'`
		out, err = runSqliteCommand(cmd)
		if err != nil {
			return err
		}
		count, _ = strconv.Atoi(strings.TrimRight(string(out), "\n"))
		if count != 1 {
			return errors.New("ERROR: install module-core-lorawan at first")
		}
	} else {
		log.Printf("lorawan doesn't exist")
	}

	cmd = `'SELECT COUNT(*) FROM networks WHERE plugin_name = "bacnetmaster"'`
	out, err = runSqliteCommand(cmd)
	if err != nil {
		return err
	}
	count, _ = strconv.Atoi(strings.TrimRight(string(out), "\n"))
	if count > 0 {
		log.Printf("bacnetmaster exist")
		cmd = `'SELECT COUNT(*) FROM plugins WHERE name = "module-core-bacnetmaster"'`
		out, err = runSqliteCommand(cmd)
		if err != nil {
			return err
		}
		count, _ = strconv.Atoi(strings.TrimRight(string(out), "\n"))
		if err != nil {
			return err
		}
		if count != 1 {
			return errors.New("ERROR: install module-core-bacnetmaster at first")
		}
	} else {
		log.Printf("bacnetmaster doesn't exist")
	}

	cmd = `'SELECT COUNT(*) FROM networks WHERE plugin_name = "modbus"'`
	out, err = runSqliteCommand(cmd)
	if err != nil {
		return err
	}
	count, _ = strconv.Atoi(strings.TrimRight(string(out), "\n"))
	if count > 0 {
		log.Printf("modbus exist")
		cmd = `'SELECT COUNT(*) FROM plugins WHERE name = "module-core-modbus"'`
		out, err = runSqliteCommand(cmd)
		if err != nil {
			return err
		}
		count, _ = strconv.Atoi(strings.TrimRight(string(out), "\n"))
		if err != nil {
			return err
		}
		if count != 1 {
			return errors.New("ERROR: install module-core-modbus at first")
		}
	} else {
		log.Printf("modbus doesn't exist")
	}

	cmd = `'UPDATE networks SET plugin_name = "module-core-loraraw", plugin_uuid = (SELECT uuid FROM plugins WHERE name = "module-core-loraraw") WHERE plugin_name = "lora";SELECT changes();'`
	out, err = runSqliteCommand(cmd)
	if err != nil {
		return err
	}
	log.Printf("Rows affected for module-core-loraraw: %s", string(out))

	cmd = `'UPDATE networks SET plugin_name = "module-core-lorawan", plugin_uuid = (SELECT uuid FROM plugins WHERE name = "module-core-lorawan") WHERE plugin_name = "lorawan";SELECT changes();'`
	out, err = runSqliteCommand(cmd)
	if err != nil {
		return err
	}
	log.Printf("Rows affected for module-core-lorawan: %s", string(out))

	cmd = `'UPDATE networks SET plugin_name = "module-core-modbus", plugin_uuid = (SELECT uuid FROM plugins WHERE name = "module-core-modbus") WHERE plugin_name = "modbus";SELECT changes();'`
	out, err = runSqliteCommand(cmd)
	if err != nil {
		return err
	}
	log.Printf("Rows affected for module-core-modbus: %s", string(out))

	cmd = `'UPDATE networks SET plugin_name = "module-core-bacnetmaster", plugin_uuid = (SELECT uuid FROM plugins WHERE name = "module-core-bacnetmaster") WHERE plugin_name = "bacnetmaster";SELECT changes();'`
	out, err = runSqliteCommand(cmd)
	if err != nil {
		return err
	}
	log.Printf("Rows affected for module-core-bacnetmaster: %s", string(out))
	return nil
}

func downloadRosDb(ip, sshUsername, sshPassword, sshPort string) error {
	cmd := fmt.Sprintf("sshpass -p '%s' scp -o StrictHostKeyChecking=no -P %s %s@%s:%s %s",
		sshPassword, sshPort, sshUsername, ip, rosDbFile, rosDownloadDbFile)
	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func uploadRosDb(ip, sshUsername, sshPassword, sshPort string) error {
	cmd := fmt.Sprintf("sshpass -p '%s' scp -o StrictHostKeyChecking=no -P %s %s %s@%s:%s",
		sshPassword, sshPort, rosDownloadDbFile, sshUsername, ip, rosDbFile)
	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func runSqliteCommand(cmd string) ([]byte, error) {
	command := exec.Command("sh", "-c", fmt.Sprintf("sqlite3 %s %s", rosDownloadDbFile, cmd))
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return output, nil
}

func restartROS(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run("sudo systemctl restart nubeio-rubix-os.service")
}
