/*
  __  __                      _  _____
 |  \/  |                    | |/ / _ \
 | \  / | __ _  ___ _ __ ___ | ' | (_) |___
 | |\/| |/ _` |/ __| '__/ _ \|  < > _ </ __|
 | |  | | (_| | (__| | | (_) | . | (_) \__ \
 |_|  |_|\__,_|\___|_|  \___/|_|\_\___/|___/

  Written by Mutlu Polatcan
  14.07.2020
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mpolatcan/macrok8s/internal"
)

func main() {
	macroK8s := internal.NewMacroK8s()

	// ==================================== CREATE COMMAND FLAGS =============================================
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createdClusterName := createCmd.String("name", "default", "Name of Microk8s's cluster")
	masterCpus := createCmd.String("master_cpus", "1", "Master instance CPU count")
	masterMem := createCmd.String("master_mem", "1", "Master instance memory size in GB")
	masterDisk := createCmd.String("master_disk", "1", "Master instance disk size in GB")
	workerCount := createCmd.Int("worker_count", 1, "Cluster worker count")
	workerCpus := createCmd.String("worker_cpus", "1", "Worker CPU count")
	workerMem := createCmd.String("worker_mem", "1G", "Worker memory size in MB or GB")
	workerDisk := createCmd.String("worker_disk", "1G", "Worker disk size in MB or GB")
	k8sVersion := createCmd.String("k8s_version", "latest/stable", "Kubernetes version")
    // ========================================================================================================

    // ======================================= DELETE COMMAND FLAGS ===========================================
    deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
    deletedClusterName := deleteCmd.String("name", "default", "Name of Microk8s's cluster")
    // ========================================================================================================

	// ======================================= START COMMAND FLAGS ===========================================
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startedClusterName := startCmd.String("name", "default", "Name of Microk8s's cluster")
	// ========================================================================================================

	// ======================================= STOP COMMAND FLAGS ===========================================
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	stoppedClusterName := stopCmd.String("name", "default", "Name of Microk8s's cluster")
	// ========================================================================================================

	switch os.Args[1] {
	case "create":
		err := createCmd.Parse(os.Args[2:])

		if err != nil {
			createCmd.Usage()
		} else {
			macroK8s.CreateCluster(*createdClusterName, *masterCpus, *masterMem,
												 *masterDisk, *workerCount, *workerCpus,
												 *workerMem, *workerDisk, *k8sVersion)
		}
	case "delete":
		err := deleteCmd.Parse(os.Args[2:])

		if err != nil {
			deleteCmd.Usage()
		} else {
			macroK8s.DeleteCluster(*deletedClusterName)
		}
	case "start": {
		err := startCmd.Parse(os.Args[2:])

		if err != nil {
			startCmd.Usage()
		} else {
			macroK8s.StartCluster(*startedClusterName)
		}
	}
	case "stop": {
		err := stopCmd.Parse(os.Args[2:])

		if err != nil {
			stopCmd.Usage()
		} else {
			macroK8s.StopCluster(*stoppedClusterName)
		}
	}
	default:
		fmt.Printf("Wrong command \"%s\"! Exitingsss...", os.Args[1])
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		println("Exiting...")
		os.Exit(1)
	}
}