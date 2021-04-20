package csharp

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	binName := "a.exe"

	// TODO: getting JSON package to work with C#
	// https://stackoverflow.com/questions/50079804/how-do-i-add-json-assembly-to-compile-c-sharp-code-on-ubuntu-terminal-using-mcs
	// https://stackoverflow.com/questions/49767212/how-to-compile-and-run-c-sharp-project-from-terminal-mac-os
	sourceFiles := util.FilterByExtension(files, "cs")
	args := append([]string{"mcs", "-out:" + binName}, sourceFiles...)
	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	binPath := filepath.Join(workDir, binName)
	return cmd.RunStdin(workDir, stdin, "mono", binPath)
}