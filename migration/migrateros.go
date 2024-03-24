package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func BackupAndMigrateROS() {
	source := "/data/rubix-os/data/data.db"
	currentDateTime := time.Now().UTC().Format("20060102150405")
	destinationDir := fmt.Sprintf("/data/backup/migration/rubix-os/%s", currentDateTime)
	destination := filepath.Join(destinationDir, "data.db")
	err := backup(source, destinationDir, destination)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	migrateROS(source)
}

func backup(source, destinationDir, destination string) error {
	_ = os.MkdirAll(destinationDir, os.FileMode(755))
	sourceFile, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()
	return err
}

func migrateROS(source string) {
	db, err := sql.Open("sqlite3", source)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM networks WHERE plugin_name = 'lora'").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		log.Printf("loraraw exist")
		err = db.QueryRow("SELECT COUNT(*) FROM plugins WHERE name = 'module-core-loraraw'").Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		if count != 1 {
			log.Fatal("ERROR: install module-core-loraraw at first")
		}
	} else {
		log.Printf("loraraw doesn't exist")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM networks WHERE plugin_name = 'lorawan'").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		log.Printf("lorawan exist")
		err = db.QueryRow("SELECT COUNT(*) FROM plugins WHERE name = 'module-core-lorawan'").Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		if count != 1 {
			log.Fatal("ERROR: install module-core-lorawan at first")
		}
	} else {
		log.Printf("lorawan doesn't exist")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM networks WHERE plugin_name = 'bacnetmaster'").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		log.Printf("bacnetmaster exist")
		err = db.QueryRow("SELECT COUNT(*) FROM plugins WHERE name = 'module-core-bacnetmaster'").Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		if count != 1 {
			log.Fatal("ERROR: install module-core-bacnetmaster at first")
		}
	} else {
		log.Printf("bacnetmaster doesn't exist")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM networks WHERE plugin_name = 'modbus'").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		log.Printf("modbus exist")
		err = db.QueryRow("SELECT COUNT(*) FROM plugins WHERE name = 'module-core-modbus';").Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		if count != 1 {
			log.Fatal("ERROR: install module-core-modbus at first")
		}
	} else {
		log.Printf("modbus doesn't exist")
	}

	result, err := db.Exec("UPDATE networks SET plugin_name = 'module-core-loraraw', plugin_uuid = (SELECT uuid FROM plugins WHERE name = 'module-core-loraraw') WHERE plugin_name = 'lora'")
	if err != nil {
		panic(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rows affected for module-core-loraraw: %d\n", rowsAffected)

	result, err = db.Exec("UPDATE networks SET plugin_name = 'module-core-lorawan', plugin_uuid = (SELECT uuid FROM plugins WHERE name = 'module-core-lorawan') WHERE plugin_name = 'lorawan'")
	if err != nil {
		panic(err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rows affected for module-core-lorawan: %d\n", rowsAffected)

	result, err = db.Exec("UPDATE networks SET plugin_name = 'module-core-bacnetmaster', plugin_uuid = (SELECT uuid FROM plugins WHERE name = 'module-core-bacnetmaster') WHERE plugin_name = 'bacnetmaster'")
	if err != nil {
		panic(err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rows affected for module-core-bacnetmaster: %d\n", rowsAffected)

	result, err = db.Exec("UPDATE networks SET plugin_name = 'module-core-modbus', plugin_uuid = (SELECT uuid FROM plugins WHERE name = 'module-core-modbus') WHERE plugin_name = 'modbus'")
	if err != nil {
		panic(err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rows affected for module-core-modbus: %d\n", rowsAffected)
}
