package scala

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])

	args := append([]string{"scalac"}, files...)
	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	return cmd.RunStdin(workDir, stdin, "scala", "Main")
}