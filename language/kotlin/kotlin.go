package kotlin

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"path/filepath"
	"strings"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	fname := filepath.Base(files[0])

	// https://discuss.kotlinlang.org/t/kotlinc-example-with-separate-files/3389/4
	stdout, stderr, err, duration := cmd.Run(workDir, "kotlinc", "*.kt", "-progressive", "-no-reflect")

	// remove java warning
	stderr = strings.Replace(stderr, "Java HotSpot(TM) 64-Bit Server VM warning: Options -Xverify:none and -noverify were deprecated in JDK 13 and will likely be removed in a future release.\n", "", -1)

	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	return cmd.RunStdin(workDir, stdin, "kotlin", className(fname))
}

func className(fname string) string {
	if len(fname) < 5 {
		return fname
	}

	ext := filepath.Ext(fname)
	name := fname[0 : len(fname)-len(ext)]
	return strings.ToUpper(string(name[0])) + name[1:] + "Kt"
}