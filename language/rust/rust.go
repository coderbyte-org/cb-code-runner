package rust

import (
	"github.com/coderbyte-org/cb-code-runner/config"
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	binName := "a.out"
	binPath := filepath.Join(workDir, binName)

	// Default: rustc -o a.out <first file>
	compileArgs := []string{"rustc", "-o", binName, files[0]}

	// All .rs files (for ENV_SOURCE_FILES)
	sourceFiles := util.FilterByExtension(files, "rs")

	// Override compile via .cbconfig "compile", supporting ENV_SOURCE_FILES
	if cfgCompile := config.ParseCbConfigField(files, "compile"); cfgCompile != nil {
		var expanded []string
		for _, part := range cfgCompile {
			if part == "ENV_SOURCE_FILES" {
				expanded = append(expanded, sourceFiles...)
			} else {
				expanded = append(expanded, part)
			}
		}
		compileArgs = expanded
	}

	stdout, stderr, err, duration := cmd.Run(workDir, compileArgs...)
	if err != nil {
		// Only fail on actual command error. Cargo often writes progress to stderr.
		return stdout, stderr, err, duration
	}

	// Default run: execute compiled binary
	runArgs := []string{binPath}

	// Override run via .cbconfig "run"
	if cfgRun := config.ParseCbConfigField(files, "run"); cfgRun != nil {
		runArgs = cfgRun
	}

	// Run (e.g. cargo test) — this will now actually execute
	stdout2, stderr2, err2, duration2 := cmd.RunStdin(workDir, stdin, runArgs...)

	// Optional: combine outputs so the user sees both build + run logs
	combinedStdout := stdout + stdout2
	combinedStderr := stderr + stderr2

	return combinedStdout, combinedStderr, err2, duration2
}