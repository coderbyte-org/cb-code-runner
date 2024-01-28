package dart

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])

	// first install based on .yaml file
	stdout, stderr, err, duration := cmd.Run(workDir, "dart", "pub", "get")
		
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	// run dart program
	return cmd.RunStdin(workDir, stdin, "dart", "run", files[0])
}