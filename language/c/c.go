package c

import (
	"github.com/coderbyte-org/cb-code-runner/config"
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
	"strings"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	binName := "a.out"
	binPath := filepath.Join(workDir, binName)

	sourceFiles := util.FilterByExtension(files, "c")

	// Default compile: clang -o a.out -lm <all .c files>
	compileArgs := append([]string{"clang", "-o", binName, "-lm"}, sourceFiles...)

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
		return stdout, stderr, err, duration
	}

	// If it's a normal C compile, you may still want to fail on stderr (warnings).
	// For now, only treat stderr as fatal when the compiler is clang/gcc.
	if stderr != "" && len(compileArgs) > 0 {
		compiler := strings.ToLower(compileArgs[0])
		if compiler == "clang" || compiler == "gcc" || compiler == "cc" {
			return stdout, stderr, err, duration
		}
	}

	// Default run: execute compiled binary
	runArgs := []string{binPath}

	// Override run via .cbconfig "run"
	if cfgRun := config.ParseCbConfigField(files, "run"); cfgRun != nil {
		runArgs = cfgRun
	}

	stdout2, stderr2, err2, duration2 := cmd.RunStdin(workDir, stdin, runArgs...)

	// Optional: combine outputs so users see compile + run logs together
	return stdout + stdout2, stderr + stderr2, err2, duration2
}