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
			log.Println(stderr)
		}

		if stdout != "" {
			log.Println(stdout)
		}

		if failureMsg != "" {
			log.Fatal(failureMsg)
		}
	} else {
		log.Println(stdout)

		if successMsg != "" {
			log.Println(successMsg)
		}
	}

	return stdout
}