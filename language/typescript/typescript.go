package typescript

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	jsName := "a.js"

	// Find all typescript files and build compile command
	sourceFiles := util.FilterByExtension(files, "ts")
	args := append([]string{"tsc", "-incremental", "-skipLibCheck", "-target", "esnext", "-declaration", "-module", "system", "-outFile", jsName}, sourceFiles...)

	// Compile to javascript
	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	if err != nil {
		return stdout, stderr, err, duration
	}

	return cmd.RunStdin(workDir, stdin, "node", jsName)
}