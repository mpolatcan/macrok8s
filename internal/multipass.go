/* Written by Mutlu Polatcan */
package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type MultipassManager struct {}

type MultipassNode struct {
	IpV4 []string `json:"ipv4"`
	Name string `json:"name"`
	Release string `json:"release"`
	State string `json:"state"`
}

type MultipassNodeList struct {
	List []MultipassNode `json:"list"`
}

func (mngr *MultipassManager) RunCmd(successMsg string, failureMsg string, args ...string) string {
	return RunCmd("multipass", successMsg, failureMsg, args...)
}

func (mngr *MultipassManager) CreateNode(name string, cpu string, mem string, disk string) {
	log.Printf("Creating node \"%s\"...", name)
	mngr.RunCmd(fmt.Sprintf("Node \"%s\" created successfully!", name),
         		fmt.Sprintf("Node \"%s\" can't created successfully!", name),
         		"launch", "--name", name, "--cpus", cpu, "--mem", mem, "--disk", disk)
}

func (mngr *MultipassManager) DeleteNode(name string) {
	log.Printf("Deleting node \"%s\"...", name)
	mngr.RunCmd(fmt.Sprintf("Node \"%s\" deleted successfully!", name),
				fmt.Sprintf("Node \"%s\" can't deleted successfully!", name),
				"delete", name)
	mngr.RunCmd("", "", "purge")
}

func (mngr *MultipassManager) StopNode(name string) {
	log.Printf("Stopping node \"%s\"...", name)
	mngr.RunCmd(fmt.Sprintf("Node \"%s\" stopped successfully!", name),
				fmt.Sprintf("Node \"%s\" can't stopped successfully!", name),
				"stop", name)
}

func (mngr *MultipassManager) StartNode(name string) {
	log.Printf("Starting node \"%s\"...", name)
	mngr.RunCmd(fmt.Sprintf("Node \"%s\" started successfully!", name),
				fmt.Sprintf("Node \"%s\" can't started successfully!", name),
				"start", name)
}

func (mngr *MultipassManager) ExecNode(name string, cmd ...string) {
	log.Printf("Executing command \"%s\" on node \"%s\"...", strings.Join(cmd, " "), name)
	mngr.RunCmd(fmt.Sprintf("Node \"%s\" executed command successfully!", name),
			    fmt.Sprintf("Node \"%s\" can't executed command succcessfully!", name),
				append([]string{"exec", name, "--"}, cmd...)...)
}

func (mngr *MultipassManager) GetNodes(clusterName string) []MultipassNode {
	nodes := mngr.RunCmd("", "", "ls", "--format", "json")

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

func (mngr *MultipassManager) Purge() {
	mngr.RunCmd("", "", "purge")
}