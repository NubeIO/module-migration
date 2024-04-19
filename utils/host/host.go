package host

import (
	"github.com/NubeIO/module-migration/cli"
	"github.com/NubeIO/module-migration/utils/file"
	"log"
	"strconv"
)

type Host struct {
	LocationUUID   string `json:"location_uuid"`
	LocationName   string `json:"location_name"`
	GroupUUID      string `json:"group_uuid"`
	GroupName      string `json:"group_name"`
	HostUUID       string `json:"host_uuid"`
	HostName       string `json:"host_name"`
	VirtualIP      string `json:"virtual_ip"`
	RosMigration   bool   `json:"ros_migration"`
	WiresMigration bool   `json:"wires_migration"`
	PluginDeletion bool   `json:"plugin_deletion"`
}

var hostFilePath = "/data/migration/migration.csv"

func GenerateHosts() {
	hosts, err := GetHosts()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	if len(hosts) > 0 {
		return
	}

	if err = createHosts(); err != nil {
		log.Printf(err.Error())
	}
}

func GetHosts() ([]*Host, error) {
	records, err := file.ReadCsvFile(hostFilePath)
	if err != nil {
		return nil, err
	}

	return mapToHosts(records), nil
}

func UpdateHosts(hosts []*Host) error {
	var data [][]string
	for _, host := range hosts {
		data = append(data, []string{
			host.LocationUUID,
			host.LocationName,
			host.GroupUUID,
			host.GroupName,
			host.HostUUID,
			host.HostName,
			host.VirtualIP,
			strconv.FormatBool(host.RosMigration),
			strconv.FormatBool(host.WiresMigration),
			strconv.FormatBool(host.PluginDeletion),
		})
	}
	return file.WriteCsvFile(hostFilePath, data)
}

func createHosts() error {
	locations, err := cli.CLI.GetLocations()
	if err != nil {
		return err
	}
	var data [][]string
	for _, loc := range locations {
		for _, grp := range loc.Groups {
			for _, hos := range grp.Hosts {
				data = append(data, []string{
					loc.UUID,
					loc.Name,
					grp.UUID,
					grp.Name,
					hos.UUID,
					hos.Name,
					hos.VirtualIP,
					"false",
					"false",
					"false",
				})
			}
		}
	}

	return file.WriteCsvFile(hostFilePath, data)
}

func mapToHosts(records [][]string) []*Host {
	hosts := make([]*Host, 0)
	for _, record := range records {
		hosts = append(hosts,
			&Host{
				LocationUUID:   record[0],
				LocationName:   record[1],
				GroupUUID:      record[2],
				GroupName:      record[3],
				HostUUID:       record[4],
				HostName:       record[5],
				VirtualIP:      record[6],
				RosMigration:   record[7] == "true",
				WiresMigration: record[8] == "true",
				PluginDeletion: record[9] == "true",
			})
	}
	return hosts
}
