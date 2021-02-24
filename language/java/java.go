package java

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	fname := filepath.Base(files[0])

	stdout, stderr, err, duration := cmd.Run(workDir, "javac", fname)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	return cmd.RunStdin(workDir, stdin, "java", className(fname))
}

func className(fname string) string {
	ext := filepath.Ext(fname)
	return fname[0 : len(fname)-len(ext)]
}