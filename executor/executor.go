package executor

import (
	"bytes"
	"errors"
	"os/exec"
)

//ExecuteCommand ...
func ExecuteCommand(command string, argv ...string) (string, error) {
	exe := exec.Command(command, argv...)

	var stdOutBuffer, stdErrBuffer bytes.Buffer

	exe.Stderr = &stdErrBuffer
	exe.Stdout = &stdOutBuffer

	if err := exe.Start(); err != nil {
		return "", errors.New(stdErrBuffer.String())
	}

	if err := exe.Wait(); err != nil {
		return "", errors.New(stdErrBuffer.String())
	}

	return stdOutBuffer.String(), nil
}
