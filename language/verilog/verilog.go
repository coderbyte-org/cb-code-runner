package verilog

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	fname := filepath.Base(files[0])
	binName := "main"

	stdout, stderr, err, duration := cmd.Run(workDir, "iverilog", "-o", binName, fname)
		
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	return cmd.RunStdin(workDir, stdin, "vvp", binName)
}