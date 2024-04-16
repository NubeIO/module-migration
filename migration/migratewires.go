package migration

import (
	"encoding/json"
	"github.com/NubeIO/flow-eng/nodes"
	"github.com/NubeIO/module-migration/utils/wiresnew"
	"github.com/NubeIO/module-migration/utils/wiresold"
	"log"
)

func MigrateWires() {
	nodeList, hostUUID, err := wiresold.Get()
	if err != nil {
		log.Fatal(err)
		return
	}

	var encodedNodes nodes.NodesList
	err = json.Unmarshal(nodeList, &encodedNodes)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = wiresnew.Migrate(&wiresnew.FlowDownload{
		HostUUID:     hostUUID,
		EncodedNodes: &encodedNodes,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
}
