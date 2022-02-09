package operations

import (
	"bytes"
	"os/exec"
)

func execute(workingDirectory, command string, arguments []string) (string, string, error) {

	cmd := exec.Command(command, arguments...)
	cmd.Dir = workingDirectory

	var stderrStr bytes.Buffer
	cmd.Stderr = &stderrStr

	var stdoutStr bytes.Buffer
	cmd.Stdout = &stdoutStr

	err := cmd.Start()

	if err != nil {
		return "", "", err
	}

	cmd.Wait()

	return stdoutStr.String(), stderrStr.String(), nil
}
