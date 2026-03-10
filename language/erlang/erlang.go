package erlang

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
	"strings"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	sourceFiles := util.FilterByExtension(files, "erl")

	// Compile all .erl files
	compileArgs := append([]string{"erlc"}, sourceFiles...)
	stdout, stderr, err, duration := cmd.Run(workDir, compileArgs...)
	if err != nil {
		return stdout, stderr, err, duration
	}

	// Get module name from first source file (without .erl extension)
	moduleName := strings.TrimSuffix(filepath.Base(sourceFiles[0]), ".erl")

	// Run the compiled module
	runArgs := []string{"erl", "-noshell", "-s", moduleName, "main", "-s", "init", "stop"}
	stdout2, stderr2, err2, duration2 := cmd.RunStdin(workDir, stdin, runArgs...)

	return stdout + stdout2, stderr + stderr2, err2, duration2
}
