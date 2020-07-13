package internal

import (
	"fmt"
	"strings"
)

type MicroK8s struct {}

func (microK8s *MicroK8s) RunCmd(args ...string) string {
	return microK8s.RunCmdWithMsg("", "", args...)
}

func (microK8s *MicroK8s) RunCmdWithMsg(successMsg string, failureMsg string, args ...string) string {
	return RunCmd("microk8s", successMsg, failureMsg, args...)
}

func (microK8s *MicroK8s) Install(cpu string, mem string, disk string, k8sVersion string) {
	microK8s.RunCmdWithMsg("MicroK8s master node installed successfully!",
		                   "MicroK8s master node can't installed successfully!",
		                   "install",
		                    "--cpu", cpu,
		                    "--mem", mem,
		                    "--disk", disk,
		                    "--channel", k8sVersion)
}

func (microK8s *MicroK8s) Start() {
	microK8s.RunCmd("start")
}

func (microK8s *MicroK8s) Stop() {
	microK8s.RunCmd("stop")
}

func (microK8s *MicroK8s) RemoveNode(nodeName string) {
	microK8s.RunCmdWithMsg(fmt.Sprintf("Node \"%s\" removed from cluster successfully!", nodeName),
		                   fmt.Sprintf("Node \"%s\" can't removed from cluster successfully!", nodeName),
		                   "remove-node", nodeName)
}

func (microK8s *MicroK8s) AddNode(nodeName string) {
	joinNodeCmdTokens := strings.Split(microK8s.RunCmd("add-node"), " ")
	joinNodeCmd := joinNodeCmdTokens[len(joinNodeCmdTokens)-6:len(joinNodeCmdTokens)-3]

	(&Multipass{}).ExecNode(nodeName, append([]string{"sudo"}, joinNodeCmd...)...)
}