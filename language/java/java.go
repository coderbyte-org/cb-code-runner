package java

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	fname := filepath.Base(files[0])

	// need to include junit into classpath env variable along with any other .jar files
	// https://github.com/junit-team/junit4/wiki/Getting-started
	jarFileIncludes := ".:/usr/local/bin/JUNIT/junit-4.13.2.jar:/usr/local/bin/JUNIT/hamcrest-core-1.3.jar:/usr/local/bin/JUNIT/gson-2.10.1.jar:/usr/local/bin/JUNIT/okhttp-3.9.1.jar:/usr/local/bin/JUNIT/okio-1.9.0.jar"
	stdout, stderr, err, duration := cmd.Run(workDir, "javac", "-Xlint:unchecked", "-cp", jarFileIncludes, fname)
		
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	return cmd.RunStdin(workDir, stdin, "java", "-cp", jarFileIncludes, className(fname))
}

func className(fname string) string {
	ext := filepath.Ext(fname)
	return fname[0 : len(fname)-len(ext)]
}