package apex

import (
	"os"
	"path/filepath"
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/config"
	"github.com/coderbyte-org/cb-code-runner/util"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])

	// Collect possible Apex source files for ENV_SOURCE_FILES
	sourceFiles := append([]string{},
		util.FilterByExtension(files, "apex")...,
	)
	sourceFiles = append(sourceFiles, util.FilterByExtension(files, "cls")...)
	sourceFiles = append(sourceFiles, util.FilterByExtension(files, "trigger")...)

	// need to auth first
	// https://developer.salesforce.com/docs/atlas.en-us.sfdx_dev.meta/sfdx_dev/sfdx_dev_auth_jwt_flow.htm
	authArgs := []string{
		"sf", "org", "login", "jwt",
		"--client-id", os.Getenv("CONSUMER_ID_SALESFORCE"),
		"--jwt-key-file", "/usr/local/bin/JWT/server.key",
		"--username", "support@coderbyte.com",
		"--alias", "my-hub-org",
		"--set-default-dev-hub",
	}

	authStdout, authStderr, authErr, authDuration := cmd.Run(workDir, authArgs...)
	if authErr != nil {
		return authStdout, authStderr, authErr, authDuration
	}

	// default run command
	runArgs := []string{
		"sf", "apex", "run",
		"--file", files[0],
		"--target-org", "my-hub-org",
	}

	// Allow .cbconfig override of "run"
	if cfgRun := config.ParseCbConfigField(files, "run"); cfgRun != nil {
		var expanded []string
		for _, part := range cfgRun {
			switch part {
			case "ENV_SOURCE_FILES":
				expanded = append(expanded, sourceFiles...)
			default:
				expanded = append(expanded, part)
			}
		}
		runArgs = expanded
	}

	runStdout, runStderr, runErr, runDuration := cmd.RunStdin(workDir, stdin, runArgs...)

	// Combine auth + run logs
	combinedStdout := authStdout + runStdout
	combinedStderr := authStderr + runStderr

	return combinedStdout, combinedStderr, runErr, runDuration
}