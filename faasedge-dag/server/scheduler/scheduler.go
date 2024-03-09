package scheduler

import (
	"faasedge-dag/server/dag"
	"log"
	"strconv"
)

var nodeInfoMap = map[string]*NodeInfo {
	"node1": {IpAddr: "192.168.1.1"},
	"node2": {IpAddr: "192.168.1.2"},
	"node3": {IpAddr: "192.168.1.3"},
	"node4": {IpAddr: "192.168.1.4"},
}

var Scheduler = BaseScheduler{}
var baseFaasWrapper = BaseFaaSWrapper{}
var vDagList = make(map[string]*VDag)

type BaseScheduler struct{}

func (bs *BaseScheduler) ScheduleDag(vDag *VDag) map[string]string {
	keys := make([]string, 0, len(nodeInfoMap))
	for k := range nodeInfoMap {
		keys = append(keys, k)
	}

	currNodeIdx := 0

	res := make(map[string]string)

	for _, function := range vDag.DagDefinition.Functions {
		nodeInfo := nodeInfoMap[keys[currNodeIdx]]
		pDagDeployment := baseFaasWrapper.DeployFunction(function.Name, nodeInfo)
		res[function.Name] = pDagDeployment.Node.IpAddr + ":"+ pDagDeployment.ContainerPort
		vDag.PDagMap[function.Name] = append(vDag.PDagMap[function.Name], pDagDeployment)
		currNodeIdx = (currNodeIdx + 1) % len(nodeInfoMap)
	}

	log.Println(vDag)
	vDagList[vDag.ClientId] = vDag
	return res
}

type BaseFaaSWrapper struct{}

func (bfw BaseFaaSWrapper) DeployFunction(functionName string, node *NodeInfo) *PDagDeployment {
	return &PDagDeployment{
		ContainterName: functionName,
		Node:           node,
		ContainerPort:  strconv.Itoa(node.GetNewPortNumber()),
	}
}


func (bfw BaseFaaSWrapper) InvokeFunction(functionName string, node *NodeInfo) {
	log.Println("Invoked function, ", functionName, "on node: ", node)
}


type NodeInfo struct {
	nextAvailablePort int
	IpAddr string
}

func (n *NodeInfo) GetNewPortNumber() int {
	var res int

	if (n.nextAvailablePort == 0) {
		n.nextAvailablePort = 3000
	}

	res = n.nextAvailablePort
	n.nextAvailablePort++
	return res
}

type VDag struct {
	ClientId string
	DagDefinition dag.Dag

	// mapping from function name to list of PDagDeployment objects.
	// we store a list because one vDag can be mapped to multiple physical deployments
	PDagMap map[string][]*PDagDeployment 
}

type PDagDeployment struct {
	ContainterName string
	Node           *NodeInfo
	ContainerPort  string
}

// Scheduler interface 
type DagScheduler interface {

	// takes a vDag object and tries to deploy it to physical nodes
	ScheduleDag(vdag *VDag) map[string]string
}

// Wrapper to interact with the cluster
type FaaSWrapper interface {
	// deploys a given function to a physical node
	DeployFunction(functionName string, node *NodeInfo) *PDagDeployment
	InvokeFunction(functionName string, node *NodeInfo)
}
