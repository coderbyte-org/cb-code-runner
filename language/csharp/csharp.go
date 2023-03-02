package csharp

import (
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
	"io"
	"os"
	"fmt"
)

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	binName := "a.exe"
	sourceFiles := util.FilterByExtension(files, "cs")

	// for C#, need to link .dll files to this project and include them in directory
	// https://stackoverflow.com/questions/50079804/how-do-i-add-json-assembly-to-compile-c-sharp-code-on-ubuntu-terminal-using-mcs
	// https://stackoverflow.com/questions/49767212/how-to-compile-and-run-c-sharp-project-from-terminal-mac-os
	// https://docs.unity3d.com/462/Documentation/Manual/UsingDLL.html
	oldLocation := "/usr/local/bin/Newtonsoft.Json.dll"
	newLocation := workDir + "/Newtonsoft.Json.dll"
	srcFile, err := os.Open(oldLocation)
	destFile, err := os.Create(newLocation)
	_, err = io.Copy(destFile, srcFile)
	err = destFile.Sync()
	
	args := append([]string{"mcs", "-r:Newtonsoft.Json.dll", "-out:" + binName}, sourceFiles...)

	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	binPath := filepath.Join(workDir, binName)
	return cmd.RunStdin(workDir, stdin, "mono", binPath)
}