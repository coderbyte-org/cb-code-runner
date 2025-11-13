package typescript

import (
	"github.com/coderbyte-org/cb-code-runner/config"
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	jsName := "main.js"

	// Find all typescript files and build compile command
	sourceFiles := util.FilterByExtension(files, "ts")

	// Default compile base: tsc -skipLibCheck -target esnext -declaration -module commonjs
	compileBase := []string{"tsc", "-skipLibCheck", "-target", "esnext", "-declaration", "-module", "commonjs"}

	// If .cbconfig has "compile", override the base compile command
	if cfgCompile := config.ParseCbConfigField(files, "compile"); cfgCompile != nil {
		compileBase = cfgCompile
	}

	// Always append the TS source files to the compile command
	compileArgs := append(compileBase, sourceFiles...)

	// Compile TypeScript to JavaScript
	stdout, stderr, err, duration := cmd.Run(workDir, compileArgs...)
	if err != nil {
		return stdout, stderr, err, duration
	}

	// Default run command: node main.js
	runArgs := []string{"node", jsName}

	// If .cbconfig has "run", override it entirely
	if cfgRun := config.ParseCbConfigField(files, "run"); cfgRun != nil {
		runArgs = cfgRun
	}

	return cmd.RunStdin(workDir, stdin, runArgs...)
}