/*
   Written by Mutlu Polatcan
   14.07.2020
*/
package internal

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func RunCmd(command string, successMsg string, failureMsg string, args ...string) string {
	cmd := exec.Command(command, args...)

	var stdoutBuff, stderrBuff bytes.Buffer
	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff

	err := cmd.Run()

	stdout := strings.ReplaceAll(stdoutBuff.String(), "\n", "")
	stderr := strings.ReplaceAll(stderrBuff.String(), "\n", "")

	if err != nil {
		if stderr != "" {
			log.Print(stderr)
		}

		if stdout != "" {
			log.Print(stdout)
		}

		if failureMsg != "" {
			log.Fatal(failureMsg)
		}
	} else {
		if successMsg != "" {
			log.Print(successMsg)
		}
	}

	return stdout
}