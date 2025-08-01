package kotlin

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
	"strings"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	fname := filepath.Base(files[0])

	// need to include junit into classpath env variable along with any other .jar files
	// https://github.com/junit-team/junit4/wiki/Getting-started
	jarFileIncludes := "/usr/local/bin/JUNIT/junit-4.13.2.jar:/usr/local/bin/JUNIT/hamcrest-core-1.3.jar:/usr/local/bin/JUNIT/gson-2.10.1.jar:/usr/local/bin/JUNIT/okio-1.9.0.jar:/usr/local/bin/JUNIT/okhttp-3.9.1.jar:/usr/local/bin/JUNIT/approvaltests-4.1.jar:."

	// files we need to compile
	// https://discuss.kotlinlang.org/t/kotlinc-example-with-separate-files/3389/4
	compileCommand := []string{"kotlinc", "-progressive", "-no-reflect", "-cp", jarFileIncludes}
	for i := 0; i < len(files); i++ {
		compileCommand = append(compileCommand, strings.ReplaceAll(files[i], workDir + "/", ""))
	}

	compileStdout, compileStderr, compileErr, duration := cmd.Run(workDir, compileCommand...)

	// only block if there's a true compilation error (i.e. kotlinc exited non-zero AND no class was produced)
	if compileErr != nil {
		return compileStdout, compileStderr, compileErr, duration
	}

	// run the program regardless of warnings in compileStderr
	runStdout, runStderr, runErr, _ := cmd.RunStdin(workDir, stdin, "kotlin", "-cp", jarFileIncludes, className(fname))

	// combine compile-time warnings with runtime stderr
	combinedStderr := strings.TrimSpace(compileStderr + "\n" + runStderr)

	return runStdout, combinedStderr, runErr, duration
}

func className(fname string) string {
	if len(fname) < 5 {
		return fname
	}

	ext := filepath.Ext(fname)
	name := fname[0 : len(fname)-len(ext)]
	return strings.ToUpper(string(name[0])) + name[1:] + "Kt"
}