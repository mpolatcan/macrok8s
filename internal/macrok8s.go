/*
   Written by Mutlu Polatcan
   14.07.2020
*/
package internal

import (
	"fmt"
	"log"
	"strings"
)

const InstallMicroK8sCmd  = "sudo snap install microk8s --classic"
const StartMicroK8sCmd = "sudo microk8s start"
const StopMicroK8sCmd  = "sudo microk8s stop"
const LeaveMicroK8sCmd  =  "sudo microk8s leave"

type MacroK8s struct {
	MicroK8s *MicroK8s
	Multipass *Multipass
}

func NewMacroK8s() *MacroK8s {
	MacroK8s := &MacroK8s{}
	MacroK8s.Multipass = &Multipass{}
	MacroK8s.MicroK8s = &MicroK8s{}

	return MacroK8s
}

func (macroK8s *MacroK8s) CreateCluster(clusterName string, masterCpus string, masterMem string,
	                                          masterDisk string, workerCount int, workerCpus string,
	                                          workerMem string, workerDisk string, k8sVersion string) {
	// Prepare MicroK8s master node
	macroK8s.MicroK8s.Install(masterCpus, masterMem, masterDisk, k8sVersion)
	macroK8s.MicroK8s.Start()

	// Prepare worker nodes
	for workerNodeIdx := 0; workerNodeIdx < workerCount; workerNodeIdx++ {
		workerNodeName := fmt.Sprintf("%s-worker-%d", clusterName, workerNodeIdx)
		isWorkerNodeExist := false

		// Check if node exist
		for _, node := range macroK8s.Multipass.GetNodes(clusterName) {
			if node.Name == workerNodeName {
				isWorkerNodeExist = true
				break
			}
		}

		if !isWorkerNodeExist {
			// Create node with specified specs
			macroK8s.Multipass.CreateNode(workerNodeName, workerCpus, workerMem, workerDisk)

			// Install MicroK8s on node
			macroK8s.Multipass.ExecNode(workerNodeName, strings.Split(InstallMicroK8sCmd, " ")...)

			// Start MicroK8S on node
			macroK8s.Multipass.ExecNode(workerNodeName, strings.Split(StartMicroK8sCmd, " ")...)

			// Add worker node to cluster
			macroK8s.MicroK8s.AddNode(workerNodeName)
		} else {
			log.Printf("Node \"%s\" already exist! Passing installation...", workerNodeName)
		}
	}
}

func (macroK8s *MacroK8s) DeleteCluster(clusterName string) {
	// Remove nodes from cluster and delete
	for _, node := range macroK8s.Multipass.GetNodes(clusterName) {
		macroK8s.MicroK8s.RemoveNode(node.IpV4[0])
		macroK8s.Multipass.ExecNode(node.Name, strings.Split(LeaveMicroK8sCmd, " ")...)
		macroK8s.Multipass.DeleteNode(node.Name)
	}
	macroK8s.Multipass.Purge()

	// Stop MicroK8s master node Kubernetes services and shutdown
	macroK8s.MicroK8s.Stop()
	macroK8s.Multipass.StopNode("microk8s-vm")
}

func (macroK8s *MacroK8s) StartCluster(clusterName string) {
	// Start master node of cluster and its Kubernetes services
	macroK8s.Multipass.StartNode("microk8s-vm")
	macroK8s.MicroK8s.Start()

	// Start worker nodes of cluster
	for _, node := range macroK8s.Multipass.GetNodes(clusterName) {
		macroK8s.Multipass.StartNode(node.Name)
		macroK8s.Multipass.ExecNode(node.Name, strings.Split(StartMicroK8sCmd, " ")...)
		macroK8s.MicroK8s.AddNode(node.Name)
	}
}

func (macroK8s *MacroK8s) StopCluster(clusterName string) {
	// Stop worker nodes of cluster
	for _, node := range macroK8s.Multipass.GetNodes(clusterName) {
		macroK8s.MicroK8s.RemoveNode(node.IpV4[0])
		macroK8s.Multipass.ExecNode(node.Name, strings.Split(LeaveMicroK8sCmd, " ")...)
		macroK8s.Multipass.ExecNode(node.Name, strings.Split(StopMicroK8sCmd, " ")...)
		macroK8s.Multipass.StopNode(node.Name)
	}

	// Stop MicroK8s master node Kubernetes services and shutdown
	macroK8s.MicroK8s.Stop()
	macroK8s.Multipass.StopNode("microk8s-vm")
}

