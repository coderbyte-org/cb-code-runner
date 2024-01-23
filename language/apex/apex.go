package apex

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])

	// https://developer.salesforce.com/docs/atlas.en-us.sfdx_dev.meta/sfdx_dev/sfdx_dev_develop_apex_run_anon.htm
	return cmd.RunStdin(workDir, stdin, "sudo", "sf", "apex", "run", "--file", files[0], "--target-org", "my-hub-org")
}