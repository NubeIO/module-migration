package host

import (
	"github.com/NubeIO/module-migration/cli"
	"github.com/NubeIO/module-migration/utils/file"
	"log"
)

type Host struct {
	LocationUUID         string `json:"location_uuid"`
	LocationName         string `json:"location_name"`
	GroupUUID            string `json:"group_uuid"`
	GroupName            string `json:"group_name"`
	HostUUID             string `json:"host_uuid"`
	HostName             string `json:"host_name"`
	IP                   string `json:"ip"`
	RosMigrationState    string `json:"ros_migration_state"`
	WiresMigrationState  string `json:"wires_migration_state"`
	PluginDeletionState  string `json:"plugin_deletion_state"`
	RosMigrationStatus   string `json:"ros_migration_status"`
	WiresMigrationStatus string `json:"wires_migration_status"`
	PluginDeletionStatus string `json:"plugin_deletion_status"`
}

var hostFilePath = "./migration.csv"

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
			host.IP,
			host.RosMigrationState,
			host.WiresMigrationState,
			host.PluginDeletionState,
			host.RosMigrationStatus,
			host.WiresMigrationStatus,
			host.PluginDeletionStatus,
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
	data = append(data, []string{
		"Location UUID",
		"Location Name",
		"Group UUID",
		"Group Name",
		"Host UUID",
		"Host Name",
		"IP",
		"ROS Migration State",
		"Wires Migration State",
		"Plugin Deletion State",
		"ROS Migration Status",
		"Wires Migration Status",
		"Plugin Deletion Status",
	})
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
					hos.IP,
					"false",
					"false",
					"false",
					"",
					"",
					"",
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
				LocationUUID:         record[0],
				LocationName:         record[1],
				GroupUUID:            record[2],
				GroupName:            record[3],
				HostUUID:             record[4],
				HostName:             record[5],
				IP:                   record[6],
				RosMigrationState:    record[7],
				WiresMigrationState:  record[8],
				PluginDeletionState:  record[9],
				RosMigrationStatus:   record[10],
				WiresMigrationStatus: record[11],
				PluginDeletionStatus: record[12],
			})
	}
	return hosts
}
