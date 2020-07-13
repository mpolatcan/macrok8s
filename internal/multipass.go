/*
   Written by Mutlu Polatcan
   14.07.2020
*/
package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Multipass struct {}

type MultipassNode struct {
	IpV4 []string `json:"ipv4"`
	Name string `json:"name"`
	Release string `json:"release"`
	State string `json:"state"`
}

type MultipassNodeList struct {
	List []MultipassNode `json:"list"`
}

func (multipass *Multipass) RunCmd(args ...string) string {
	return multipass.RunCmdWithMsg("", "", args...)
}

func (multipass *Multipass) RunCmdWithMsg(successMsg string, failureMsg string, args ...string) string {
	return RunCmd("multipass", successMsg, failureMsg, args...)
}

func (multipass *Multipass) CreateNode(nodeName string, cpu string, mem string, disk string) {
	log.Printf("Creating node \"%s\"...", nodeName)

	multipass.RunCmdWithMsg(fmt.Sprintf("Node \"%s\" created successfully!", nodeName),
		                    fmt.Sprintf("Node \"%s\" can't created successfully!", nodeName),
		                   "launch", "--name", nodeName, "--cpus", cpu, "--mem", mem, "--disk", disk)
}

func (multipass *Multipass) DeleteNode(nodeName string) {
	log.Printf("Deleting node \"%s\"...", nodeName)

	multipass.RunCmdWithMsg(fmt.Sprintf("Node \"%s\" deleted successfully!", nodeName),
		                    fmt.Sprintf("Node \"%s\" can't deleted successfully!", nodeName),
		                   "delete", nodeName)
}

func (multipass *Multipass) StopNode(nodeName string) {
	log.Printf("Stopping node \"%s\"...", nodeName)

	multipass.RunCmdWithMsg(fmt.Sprintf("Node \"%s\" stopped successfully!", nodeName),
		                    fmt.Sprintf("Node \"%s\" can't stopped successfully!", nodeName),
						   "stop", nodeName)
}

func (multipass *Multipass) StartNode(nodeName string) {
	log.Printf("Starting node \"%s\"...", nodeName)

	multipass.RunCmdWithMsg(fmt.Sprintf("Node \"%s\" started successfully!", nodeName),
		                    fmt.Sprintf("Node \"%s\" can't started successfully!", nodeName),
					        "start", nodeName)
}

func (multipass *Multipass) ExecNode(nodeName string, cmd ...string) {
	log.Printf("Executing command \"%s\" on node \"%s\"...", strings.Join(cmd, " "), nodeName)

	multipass.RunCmdWithMsg(fmt.Sprintf("Node \"%s\" executed command successfully!", nodeName),
		                    fmt.Sprintf("Node \"%s\" can't executed command succcessfully!", nodeName),
		                    append([]string{"exec", nodeName, "--"}, cmd...)...)
}

func (multipass *Multipass) GetNodes(clusterName string) []MultipassNode {
	nodes := multipass.RunCmd( "ls", "--format", "json")

	var multipassNodeList MultipassNodeList
	err := json.Unmarshal([]byte(nodes), &multipassNodeList)

	if err != nil {
		panic(err.Error())
	}

	var clusterNodes []MultipassNode

	for _, node := range multipassNodeList.List {
		if strings.HasPrefix(node.Name, clusterName) {
			clusterNodes = append(clusterNodes, node)
		}
	}

	return clusterNodes
}

func (multipass *Multipass) Purge() {
	multipass.RunCmd("purge")
}