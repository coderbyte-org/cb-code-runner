package rust

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	binName := "a.out"

	stdout, stderr, err, duration := cmd.Run(workDir, "rustc", "-o", binName, files[0])
	if err != nil {
		return stdout, stderr, err, duration
	}

	binPath := filepath.Join(workDir, binName)
	return cmd.RunStdin(workDir, stdin, binPath)
}