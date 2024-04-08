package commander

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func (bt *BuildTool) handleRunBash(params interface{}) (interface{}, error) {
	// Assert that params is a string
	cmdString, ok := params.(string)
	if !ok {
		return nil, errors.New("invalid params type for runBash command")
	}

	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run bash command: %v", err)
	}

	return nil, nil
}
