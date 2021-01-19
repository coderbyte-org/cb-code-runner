package kotlin

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
	"strings"
)

func Run(files []string, stdin string) (string, string, error) {
	workDir := filepath.Dir(files[0])
	fname := filepath.Base(files[0])

	stdout, stderr, err := cmd.RunStdin(workDir, stdin, "kotlinc", "-script", fname)

	// remove java warning
	stderr = strings.Replace(stderr, "Java HotSpot(TM) 64-Bit Server VM warning: Options -Xverify:none and -noverify were deprecated in JDK 13 and will likely be removed in a future release.\n", "", -1)

	return stdout, stderr, err
}