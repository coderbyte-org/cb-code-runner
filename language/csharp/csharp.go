package csharp

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	projFile := "Main.csproj"

	// run using `dotnet` now instead of the previous `mcs` method, this allows us to add .csproj files
	// https://stackoverflow.com/a/64646610
	args := append([]string{"dotnet", "build"}, projFile)

	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	binPath := filepath.Join(workDir, projFile)
	return cmd.RunStdin(workDir, stdin, "dotnet", "run", binPath)
}