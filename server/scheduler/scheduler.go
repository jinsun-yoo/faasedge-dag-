package scheduler

import (
	"faasedge-dag/server/dag"
)

var nodeInfoMap = map[string]*NodeInfo {
	"node1": {IpAddr: "192.168.1.1"},
	"node2": {IpAddr: "192.168.1.2"},
	"node3": {IpAddr: "192.168.1.3"},
	"node4": {IpAddr: "192.168.1.4"},
}

var Scheduler = BaseScheduler{}
var baseFaasWrapper = BaseFaaSWrapper{}

type BaseScheduler struct{}

func (bs *BaseScheduler) ScheduleDag(vDag *VDag) {
	keys := make([]string, 0, len(nodeInfoMap))
	for k := range nodeInfoMap {
		keys = append(keys, k)
	}

	currNode := 0

	for _, function := range vDag.DagDefinition.Functions {
		nodeInfo := nodeInfoMap[keys[currNode]]
		pDagDeployment := baseFaasWrapper.DeployFunction(function.Name, nodeInfo)
		vDag.PDagMap[function.Name] = append(vDag.PDagMap[function.Name], pDagDeployment)
		currNode = (currNode + 1) % len(nodeInfoMap)
	}
}

type BaseFaaSWrapper struct{}

func (bfw BaseFaaSWrapper) DeployFunction(functionName string, node *NodeInfo) *PDagDeployment {
	return &PDagDeployment{
		ContainterName: functionName,
		Node:           node,
		ContainerPort:  "3000",
	}
}

type NodeInfo struct {
	IpAddr string
}

type VDag struct {
	ClientId string

	// Definition of the Dag, defined by the user only once at the very beginning. This is the definition we added last week.
	DagDefinition dag.Dag

	// Map from function name to a list of pDagDeployment objects.
	// The reason we have map[string][]*pDagDeployment instead of map[string]*pDagDeployment is because one vDAG could be mapped to multiple pDAGs.
	PDagMap map[string][]*PDagDeployment
}

type PDagDeployment struct {
	ContainterName string
	Node           *NodeInfo
	ContainerPort  string
}

type DagScheduler interface {
	// For Each function in the vDAG, will assign a physical container to the function.
	ScheduleDag(vdag *VDag)
}

type FaaSWrapper interface {
	// Tells the node in node to deploy a container for this function.
	DeployFunction(functionName string, node *NodeInfo) *PDagDeployment
}
