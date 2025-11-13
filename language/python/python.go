package python

import (
	"github.com/coderbyte-org/cb-code-runner/config"
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	runArgs := []string{"python", files[0]}

	if cfgArgs := config.ParseCbConfigField(files, "run"); cfgArgs != nil {
		runArgs = cfgArgs
	}

	return cmd.RunStdin(workDir, stdin, runArgs...)
}