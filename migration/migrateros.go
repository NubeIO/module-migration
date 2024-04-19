package migration

import (
	"errors"
	"fmt"
	"github.com/NubeIO/module-migration/utils/sshclient"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/ssh"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

var rosDbFile = "/data/rubix-os/data/data.db"

func BackupAndMigrateROS(ip, sshUsername, sshPassword string) error {
	client, err := sshclient.New(ip, sshUsername, sshPassword)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer client.Close()

	currentDateTime := time.Now().UTC().Format("20060102150405")
	destinationDir := fmt.Sprintf("/data/backup/migration/rubix-os/%s", currentDateTime)
	destination := filepath.Join(destinationDir, "data.db")
	if err = backupROS(client, destination, destinationDir); err != nil {
		log.Printf(err.Error())
		return err
	}

	if err = migrateROS(client); err != nil {
		log.Printf(err.Error())
		return err
	}

	return restartROS(client)
}

func backupROS(client *ssh.Client, destination, destinationDir string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("mkdir -p %s && cp %s %s", destinationDir, rosDbFile, destination)
	return session.Run(cmd)
}

func migrateROS(client *ssh.Client) error {
	cmd := `'SELECT COUNT(*) FROM networks WHERE plugin_name = "lora"'`
	out, err := runSqliteCommand(client, cmd)
	if err != nil {
		return err
	}
	count, _ := strconv.Atoi(string(out))
	if count > 0 {
		log.Printf("loraraw exist")
		cmd = `'SELECT COUNT(*) FROM plugins WHERE name = "module-core-loraraw"'`
		out, err = runSqliteCommand(client, cmd)
		if err != nil {
			return err
		}
		count, _ = strconv.Atoi(string(out))
		if count != 1 {
			return errors.New("ERROR: install module-core-loraraw at first")
		}
	} else {
		log.Printf("loraraw doesn't exist")
	}

	cmd = `'SELECT COUNT(*) FROM networks WHERE plugin_name = "lorawan"'`
	out, err = runSqliteCommand(client, cmd)
	if err != nil {
		return err
	}
	count, _ = strconv.Atoi(string(out))
	if count > 0 {
		log.Printf("lorawan exist")
		cmd = `'SELECT COUNT(*) FROM plugins WHERE name = "module-core-lorawan"'`
		out, err = runSqliteCommand(client, cmd)
		if err != nil {
			return err
		}
		count, _ = strconv.Atoi(string(out))
		if count != 1 {
			return errors.New("ERROR: install module-core-lorawan at first")
		}
	} else {
		log.Printf("lorawan doesn't exist")
	}

	cmd = `'SELECT COUNT(*) FROM networks WHERE plugin_name = "bacnetmaster"'`
	if err != nil {
		return err
	}
	count, _ = strconv.Atoi(string(out))
	if count > 0 {
		log.Printf("bacnetmaster exist")
		cmd = `'SELECT COUNT(*) FROM plugins WHERE name = "module-core-bacnetmaster"'`
		out, err = runSqliteCommand(client, cmd)
		if err != nil {
			return err
		}
		count, _ = strconv.Atoi(string(out))
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
	if err != nil {
		return err
	}
	count, _ = strconv.Atoi(string(out))
	if count > 0 {
		log.Printf("modbus exist")
		cmd = `'SELECT COUNT(*) FROM plugins WHERE name = "module-core-modbus"'`
		out, err = runSqliteCommand(client, cmd)
		if err != nil {
			return err
		}
		count, _ = strconv.Atoi(string(out))
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
	out, err = runSqliteCommand(client, cmd)
	if err != nil {
		return err
	}
	log.Printf("Rows affected for module-core-loraraw: %s\n", string(out))

	cmd = `'UPDATE networks SET plugin_name = "module-core-lorawan", plugin_uuid = (SELECT uuid FROM plugins WHERE name = "module-core-lorawan") WHERE plugin_name = "lorawan";SELECT changes();'`
	out, err = runSqliteCommand(client, cmd)
	if err != nil {
		return err
	}
	log.Printf("Rows affected for module-core-lorawan: %s\n", string(out))

	cmd = `'UPDATE networks SET plugin_name = "module-core-modbus", plugin_uuid = (SELECT uuid FROM plugins WHERE name = "module-core-modbus") WHERE plugin_name = "modbus";SELECT changes();'`
	out, err = runSqliteCommand(client, cmd)
	if err != nil {
		return err
	}
	log.Printf("Rows affected for module-core-modbus: %s\n", string(out))
	return nil
}

func runSqliteCommand(client *ssh.Client, cmd string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session.CombinedOutput(fmt.Sprintf("sqlite3 %s %s", rosDbFile, cmd))
}

func restartROS(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run("sudo systemctl restart nubeio-rubix-os.service")
}
