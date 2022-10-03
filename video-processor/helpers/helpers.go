package helpers

import (
	"os/exec"
)

func RunCmd(cmd *exec.Cmd) (string, error) {
	b, err := cmd.CombinedOutput()
	if err != nil {
		return string(b), err
	}

	return string(b), nil
}
