package cpp

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	binName := "a.out"

	sourceFiles := util.FilterByExtension(files, "cpp")
	args := append([]string{"clang++", "-std=c++20", "-g", "-o", binName}, sourceFiles...)
	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	// compilation succeeded, now run binary
	binPath := filepath.Join(workDir, binName)
	stdout, stderr, err, duration = cmd.RunStdin(workDir, stdin, binPath)

	return stdout, stderr, err, duration
}