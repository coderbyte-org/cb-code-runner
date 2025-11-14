package csharp

import (
	"github.com/coderbyte-org/cb-code-runner/config"
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
	"strings"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	projFile := "Main.csproj"
	lastFile := files[len(files) - 1]

	// Default compile command
	compileArgs := []string{"dotnet", "build", projFile}

	// If .cbconfig has "compile", override the compile command
	if cfgCompile := config.ParseCbConfigField(files, "compile"); cfgCompile != nil {
		compileArgs = cfgCompile
	}

	stdout, stderr, err, duration := cmd.Run(workDir, compileArgs...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	// Special case: dotnet tests (unchanged, still uses built-in command)
	if strings.Contains(lastFile, "dotnet_test") {
		return cmd.RunStdin(workDir, stdin, "dotnet", "test", projFile, "--logger", "\"console;verbosity=detailed\"")
	}

	// Default run command
	runArgs := []string{"dotnet", "run"}

	// If .cbconfig has "run", override the run command
	if cfgRun := config.ParseCbConfigField(files, "run"); cfgRun != nil {
		runArgs = cfgRun
	}

	return cmd.RunStdin(workDir, stdin, runArgs...)
}