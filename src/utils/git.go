package utils

import (
	"os"
	"os/exec"
)

func GitClone(url, path string) error {
	cmd := exec.Command("git", "clone", url, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
