package php

import (
	"github.com/danielborowski/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error) {
	workDir := filepath.Dir(files[0])
	return cmd.RunStdin(workDir, stdin, "php", files[0])
}