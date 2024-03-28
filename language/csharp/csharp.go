package csharp

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
	"strings"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	projFile := "Main.csproj"
	lastFile := files[len(files) - 1]

	// run using `dotnet` now instead of the previous `mcs` method, this allows us to add .csproj files
	// https://stackoverflow.com/a/64646610
	args := append([]string{"dotnet", "build"}, projFile)

	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	binPath := filepath.Join(workDir, projFile)
	
	if (strings.Contains(lastFile, "dotnet_test")) {
		return cmd.RunStdin(workDir, stdin, "dotnet", "test", projFile, "--logger", "\"console;verbosity=detailed\"")
	} else {
		return cmd.RunStdin(workDir, stdin, "dotnet", "run", binPath)
	}
}