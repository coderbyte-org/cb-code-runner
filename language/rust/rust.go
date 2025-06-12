package rust

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
	"strings"
)

func contains(s []string, e string) bool {
	for _, a := range s {
			if strings.Contains(a, e) {
					return true
			}
	}
	return false
}

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	binName := "a.out"
	containsToml := contains(files, "Cargo.toml")

	// if contains toml, we run unit tests
	if (containsToml) {
		stdout, stderr, err, duration := cmd.Run(workDir, "cargo", "build")
		if err != nil {
			return stdout, stderr, err, duration
		}
		return cmd.RunStdin(workDir, stdin, "cargo", "test")
	} else {
		stdout, stderr, err, duration := cmd.Run(workDir, "rustc", "-o", binName, files[0])
		if err != nil || stderr != "" {
			return stdout, stderr, err, duration
		}
		binPath := filepath.Join(workDir, binName)
		return cmd.RunStdin(workDir, stdin, binPath)
	}
}