package wiresnew

import (
	flowctrl "github.com/NubeIO/flow-eng"
	"github.com/NubeIO/flow-eng/node"
	"github.com/NubeIO/flow-eng/nodes"
	"github.com/NubeIO/rubix-edge-wires/db"
	"github.com/mitchellh/mapstructure"
)

var dbFile = "/data/rubix-edge-wires/data/data.json"

var (
	storage  db.DB
	flowInst *flowctrl.Flow
)

type FlowDownload struct {
	HostUUID     string           `json:"hostUUID"`
	EncodedNodes *nodes.NodesList `json:"encodedNodes"`
}

func Migrate(flowDownload *FlowDownload) error {
	if flowDownload == nil {
		return nil
	}

	storage, _ = db.New(dbFile)
	flowInst = flowctrl.New()
	nodeList := &nodes.NodesList{}

	err := mapstructure.Decode(flowDownload.EncodedNodes, &nodeList)
	if err != nil {
		return err
	}

	dec, err := decode(nodeList)
	if err != nil {
		return err
	}

	return setLatestFlow(dec, true, flowDownload.HostUUID)
}

func decode(encodedNodes *nodes.NodesList) ([]*node.Spec, error) {
	return nodes.Decode(encodedNodes)
}

func setLatestFlow(flow []*node.Spec, saveFlowToDB bool, hostUUID string) error {
	if saveFlowToDB {
		saveFlowDB(flow, hostUUID)
	}
	return nil
}

func saveFlowDB(flow []*node.Spec, hostUUID string) *db.Backup {
	back := &db.Backup{Data: flow, HostUUID: hostUUID}
	bu := storage.AddBackup(back, 5)
	storage.Save()
	return bu
}
