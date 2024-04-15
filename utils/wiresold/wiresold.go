package wiresold

import (
	"encoding/json"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/helpers/store"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	bacnetio "github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"strconv"
)

var (
	dbFile       = "/data/rubix-edge-wires/data/data.db"
	latestBackup *db.Backup
	latestFlow   []*node.Spec
	flowInst     *flowctrl.Flow
	storage      db.DB
	cacheStore   *store.Store
)

func Get() (nodeList []byte, hostUUID string, err error) {
	storage = db.New(dbFile)
	if storage == nil {
		return
	}

	err = getLatestFlow()
	if err != nil {
		return
	}

	flowInst = flowctrl.New()

	err = addDefaultConnection()
	if err != nil {
		return
	}

	start()

	encode, err := nodes.Encode(flowInst.Get())
	if err != nil {
		return
	}

	nodeList, err = json.Marshal(encode)
	hostUUID = latestBackup.HostUUID

	return
}

func addDefaultConnection() error {
	c, err := storage.GetConnections()
	if err != nil {
		return err
	}
	var flowNetworkConnection = names.FlowFramework
	var found bool
	for _, connection := range c {
		if connection.Application == flowNetworkConnection {
			found = true
		}
	}
	if !found {
		_, err = storage.AddConnection(&db.Connection{
			Application: names.FlowFramework,
			Name:        "flow framework integration over MQTT (dont edit/delete)",
			Host:        "127.0.0.1",
			Port:        1883,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func start() {
	var err error
	var nodesList []node.Node
	var parentList = nodes.FilterNodes(latestFlow, nodes.FilterIsParent, "")
	var parentChildList = nodes.FilterNodes(latestFlow, nodes.FilterIsParentChild, "")
	var childList = nodes.FilterNodes(latestFlow, nodes.FilterIsChild, "")
	var nonChildNodes = nodes.FilterNodes(latestFlow, nodes.FilterNonContainer, "")

	if cacheStore == nil {
		cacheStore = makeStore()
	}

	var bacnetStore *bacnetio.Bacnet
	if bacnetStore == nil {
		app := names.Modbus
		deviceCount := "0"
		for _, n := range latestFlow {
			if n.GetName() == "bacnet-server" {
				schema, err := bacnetio.GetBacnetSchema(n.Settings)
				if err != nil {
				}
				if schema != nil {
					deviceCount = schema.DeviceCount
				}
			}
			bacnetStore = makeBacnetStore(string(app), deviceCount)
		}
	}
	var networksPool driver.Driver // flow-framework networks inst
	if networksPool == nil {
		networksPool = driver.New(&driver.Networks{})
	}

	// add the container nodes first, then the children and so on
	for _, n := range parentList {
		var node_ node.Node
		if n.Info.Category == "bacnet" {
			node_, err = nodes.Builder(n, storage, cacheStore, bacnetStore)
		} else if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, cacheStore, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage, cacheStore)
		}
		nodesList = append(nodesList, node_)
	}

	for _, n := range parentChildList {
		var node_ node.Node
		if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, cacheStore, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage, cacheStore)
		}
		nodesList = append(nodesList, node_)
	}

	for _, n := range childList {
		var node_ node.Node
		if n.Info.Category == "bacnet" {
			node_, err = nodes.Builder(n, storage, cacheStore, bacnetStore)
		} else if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, cacheStore, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage, cacheStore)
		}
		nodesList = append(nodesList, node_)
	}
	for _, n := range nonChildNodes {
		var node_ node.Node
		node_, err = nodes.Builder(n, storage, cacheStore)
		nodesList = append(nodesList, node_)
	}

	if err != nil {

	}
	flowInst.AddNodes(nodesList...)
	flowInst.MakeNodeConnections(true)
	flowInst.MakeGraph()
	for _, n := range flowInst.Get().GetNodes() { // add all nubeDevNodes to each node so data can be passed between nubeDevNodes easy
		n.AddNodes(flowInst.Get().GetNodes())
	}
}

func getLatestFlow() error {
	bac, err := storage.GetLatestBackup()
	if err != nil {
		return err
	}
	latestBackup = bac
	var nodeList []*node.Spec
	b, err := json.Marshal(bac.Data)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &nodeList); err != nil {
		return err
	}
	latestFlow = nodeList
	return nil
}

func makeStore() *store.Store {
	return store.Init()
}

func makeBacnetStore(application string, deviceCount string) *bacnetio.Bacnet {
	ip := "0.0.0.0"
	mqttClient, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{fmt.Sprintf("tcp://%s:1883", ip)},
	})
	err = mqttClient.Connect()
	if err != nil {
	}
	i, err := strconv.Atoi(deviceCount)
	app := names.ApplicationName(application)
	opts := &bacnetio.Bacnet{
		Store:       points.New(names.ApplicationName(application), nil, i, 200, 200),
		MqttClient:  mqttClient,
		Application: app,
		Ip:          ip,
	}
	return opts
}
