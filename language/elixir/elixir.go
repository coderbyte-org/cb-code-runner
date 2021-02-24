package elixir

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	sourceFiles := util.FilterByExtension(files, "ex")
	args := append([]string{"elixirc"}, sourceFiles...)
	return cmd.RunStdin(workDir, stdin, args...)
}