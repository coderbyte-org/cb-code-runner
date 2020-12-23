package cpp

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error) {
	workDir := filepath.Dir(files[0])
	binName := "a.out"

	sourceFiles := util.FilterByExtension(files, "cpp")
	args := append([]string{"clang++", "-std=c++20", "-o", binName}, sourceFiles...)
	stdout, stderr, err := cmd.Run(workDir, args...)
	if err != nil || stderr != "" {
		return stdout, stderr, err
	}

	binPath := filepath.Join(workDir, binName)
	return cmd.RunStdin(workDir, stdin, binPath)
}