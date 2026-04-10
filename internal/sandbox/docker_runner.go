package sandbox

import (
	"os/exec"
)

func RunContainer(cmd string) ([]byte, error) {
	command := exec.Command("docker", "run", "--rm", "executor-image", "bash", "-c", cmd)
	return command.CombinedOutput()
}