package wiresnew

import (
	flowctrl "github.com/NubeIO/flow-eng"
	"github.com/NubeIO/flow-eng/node"
	"github.com/NubeIO/flow-eng/nodes"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-edge-wires/db"
	"github.com/mitchellh/mapstructure"
	"strings"
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

	err = setLatestFlow(dec, true, flowDownload.HostUUID)
	if err != nil {
		return err
	}

	content, err := fileutils.ReadFile(dbFile)
	if err != nil {
		return err
	}
	newContent := replacePluginToModule(content)
	return fileutils.WriteFile(dbFile, newContent, 0644)
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

func replacePluginToModule(body string) string {
  body = strings.ReplaceAll(body, `"point":"lora:`, `"point":"module-core-loraraw:`)
  body = strings.ReplaceAll(body, `"point":"lorawan:`, `"point":"module-core-lorawan:`)
  body = strings.ReplaceAll(body, `"point":"bacnetmaster:`, `"point":"module-core-bacnetmaster:`)
  body = strings.ReplaceAll(body, `"point":"modbus:`, `"point":"module-core-modbus:`)

	return body
}
