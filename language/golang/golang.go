package golang

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	return cmd.RunStdin(workDir, stdin, "go", "run", files[0])
}