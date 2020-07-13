/* Written by Mutlu Polatcan */
package internal

import (
	"fmt"
	"strings"
)

const InstallMicroK8sCmd  = "sudo snap install microk8s --classic"
const StartMicroK8sCmd = "sudo microk8s start"

type MicroK8sClusterManager struct {
	MultipassManager *MultipassManager
}

type MultipassCmd func(name string)

func NewMicroK8sClusterManager() *MicroK8sClusterManager {
	microK8sClusterManager := &MicroK8sClusterManager{}
	microK8sClusterManager.MultipassManager = &MultipassManager{}

	return microK8sClusterManager
}

func (mngr *MicroK8sClusterManager) RunCmd(successMsg string, failureMsg string, args ...string) string {
	return RunCmd("microk8s", successMsg, failureMsg, args...)
}

func (mngr *MicroK8sClusterManager) CreateCluster(clusterName string, masterCpus string, masterMem string,
	                                              masterDisk string, workerCount int, workerCpus string,
	                                              workerMem string, workerDisk string, k8sVersion string) {
	// Prepare MicroK8s master node
	mngr.RunCmd("MicroK8s master node installed successfully!",
		        "MicroK8s master node can't installed successfully!",
		        "install",
		        "--cpu", masterCpus,
		        "--mem", masterMem,
		        "--disk", masterDisk,
		        "--channel", k8sVersion)

	// Prepare worker nodes
	for workerIdx := 0; workerIdx < workerCount; workerIdx++ {
		workerNodeName := fmt.Sprintf("%s-worker-%d", clusterName, workerIdx)
		isWorkerNodeExist := false

		// Check if node exist
		for _, node := range mngr.MultipassManager.GetNodes(clusterName) {
			if node.Name == workerNodeName {
				isWorkerNodeExist = true
				break
			}
		}

		if !isWorkerNodeExist {
			// Create node with specified specs
			mngr.MultipassManager.CreateNode(workerNodeName, workerCpus, workerMem, workerDisk)

			// Install MicroK8s on node
			mngr.MultipassManager.ExecNode(workerNodeName, strings.Split(InstallMicroK8sCmd, " ")...)

			// Start MicroK8S on node
			mngr.MultipassManager.ExecNode(workerNodeName, strings.Split(StartMicroK8sCmd, " ")...)

			// Join MicroK8S worker node to master node
			joinNodeCmdTokens := strings.Split(mngr.RunCmd("", "", "add-node"), " ")
			joinNodeCmd := joinNodeCmdTokens[len(joinNodeCmdTokens)-6:len(joinNodeCmdTokens)-3]
			mngr.MultipassManager.ExecNode(workerNodeName, append([]string{"sudo"}, joinNodeCmd...)...)
		}
	}
}

func (mngr *MicroK8sClusterManager) RunCmdOnWorkerNodes(clusterName string, cmd MultipassCmd) {
	for _, node := range mngr.MultipassManager.GetNodes(clusterName) { cmd(node.Name) }
}

func (mngr *MicroK8sClusterManager) DeleteCluster(clusterName string) {
	for _, node := range mngr.MultipassManager.GetNodes(clusterName) {
		mngr.RunCmd(fmt.Sprintf("Node \"%s\" removed from cluster successfully!", node.Name),
					fmt.Sprintf("Node \"%s\" can't removed from cluster successfully!", node.Name),
					"remove-node",
					node.IpV4[0])
	}

	mngr.RunCmdOnWorkerNodes(clusterName, mngr.MultipassManager.DeleteNode)

	mngr.MultipassManager.Purge()
}

func (mngr *MicroK8sClusterManager) StartCluster(clusterName string) {
	// Start master node of cluster
	mngr.MultipassManager.StartNode("microk8s-vm")

	// Start Kubernetes services of master node
	mngr.RunCmd("", "", "start")

	// Start worker nodes of cluster
	mngr.RunCmdOnWorkerNodes(clusterName, mngr.MultipassManager.StartNode)
}

func (mngr *MicroK8sClusterManager) StopCluster(clusterName string) {
	mngr.RunCmdOnWorkerNodes(clusterName, mngr.MultipassManager.StopNode)
	mngr.RunCmd("", "", "stop")
	mngr.MultipassManager.StopNode("microk8s-vm")
}

