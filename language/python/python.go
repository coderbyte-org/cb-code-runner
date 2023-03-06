package python

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	stdout, stderr, err, duration = cmd.RunStdin(workDir, stdin, "python", files[0])

	// remove pyspark warnings
	stdout = strings.Replace(stdout, "WARN NativeCodeLoader: Unable to load native-hadoop library for your platform... using builtin-java classes where applicable\n", "", -1)
	stderr = strings.Replace(stderr, "Setting default log level to \"WARN\".\nTo adjust logging level use sc.setLogLevel(newLevel). For SparkR, use setLogLevel(newLevel).\n", "", -1)

	return stdout, stderr, err, duration
}