package swift

import (
	"path/filepath"
	"strings"
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/config"
	"github.com/coderbyte-org/cb-code-runner/util"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])

	// Swift honors a raw shell run command from .cbconfig.
	if runCmd := strings.TrimSpace(config.ParseCbConfigFieldRaw(files, "run")); runCmd != "" {
		if stdin != "" {
			return cmd.RunStdin(workDir, stdin, "sh", "-lc", runCmd)
		}
		return cmd.Run(workDir, "sh", "-lc", runCmd)
	}

	// Default compile + run behavior.
	binName := "a.out"
	sourceFiles := util.FilterByExtension(files, "swift")
	args := append([]string{"swiftc", "-swift-version", "6", "-suppress-warnings", "-o", binName}, sourceFiles...)

	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	return cmd.RunStdin(workDir, stdin, filepath.Join(workDir, binName))
}