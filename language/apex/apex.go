package apex

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
	"os"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])

	// need to auth first
	// https://developer.salesforce.com/docs/atlas.en-us.sfdx_dev.meta/sfdx_dev/sfdx_dev_auth_jwt_flow.htm
	stdout, stderr, err, duration := cmd.Run(workDir, "sf", "org", "login", "jwt", "--client-id", os.Getenv("CONSUMER_ID_SALESFORCE"), "--jwt-key-file", "/usr/local/bin/JWT/server.key", "--username", "support@coderbyte.com", "--alias", "my-hub-org", "--set-default-dev-hub")
		
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	// run apex code
	// https://developer.salesforce.com/docs/atlas.en-us.sfdx_dev.meta/sfdx_dev/sfdx_dev_develop_apex_run_anon.htm
	return cmd.RunStdin(workDir, stdin, "sf", "apex", "run", "--file", files[0], "--target-org", "my-hub-org")
}