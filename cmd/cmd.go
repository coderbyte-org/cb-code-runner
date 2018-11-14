package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"log"
	"time"
	"errors"
)

func Run(workDir string, args ...string) (string, string, error) {
	return RunStdin(workDir, "", args...)
}

func RunStdin(workDir, stdin string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var timeoutReached bool = false

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workDir
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()

	// error starting process
	if err != nil {
		log.Println("1. process fail\n", err.Error())
	}

	// do not exceed timeout
	// https://stackoverflow.com/questions/11886531/terminating-a-process-started-with-os-exec-in-golang
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(5 * time.Second):
		err := cmd.Process.Kill()
		if err != nil {
			log.Println("2. failed to kill process\n", err.Error())
		} 
		timeoutReached = true
		log.Println("ERROR: process killed as timeout reached")
	case err := <-done:
		if err != nil {
			log.Println("3. process finished with error\n", err.Error())
		}
	}

	if timeoutReached {
		err = errors.New("process killed as timeout reached")
	}

	return stdout.String(), stderr.String(), err
}

func RunBash(workDir, command string) (string, string, error) {
	return Run(workDir, "bash", "-c", command)
}

func RunBashStdin(workDir, command, stdin string) (string, string, error) {
	return RunStdin(workDir, stdin, "bash", "-c", command)
}
