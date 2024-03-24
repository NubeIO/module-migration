package migration

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var BackupWiresPath = "/data/backup/migration/rubix-edge-wires/backup.json"

func BackupWires() {
	url := "http://localhost:1665/api/flows"

	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error while sending GET request: %s\n", err)
		return
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %s\n", err)
		return
	}
	log.Printf("Response Body: %s\n\n", string(responseBody))

	responseBody = replacePluginToModule(responseBody)
	log.Println("Migrated Body:")
	fmt.Println(string(responseBody))
	fmt.Print("\n")
	_ = os.Mkdir(filepath.Dir(BackupWiresPath), 0644)
	err = ioutil.WriteFile(BackupWiresPath, responseBody, 0644)
	if err != nil {
		log.Printf("Error saving response body to file: %s\n", err)
		return
	}
	log.Printf("Successfully migrated at: %s", BackupWiresPath)
}

func replacePluginToModule(body []byte) []byte {
	bodyString := string(body)

	bodyString = strings.ReplaceAll(bodyString, `"point":"lora`, `"point":"module-core-loraraw`)
	bodyString = strings.ReplaceAll(bodyString, `"point":"lorawan`, `"point":"module-core-lorawan`)
	bodyString = strings.ReplaceAll(bodyString, `"point":"bacnetmaster`, `"point":"module-core-bacnetmaster`)
	bodyString = strings.ReplaceAll(bodyString, `"point":"modbus`, `"point":"module-core-modbus`)

	return []byte(bodyString)
}
