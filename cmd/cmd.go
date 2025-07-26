package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"log"
	"time"
	"errors"
	"strconv"
	"syscall"
)

func Run(workDir string, args ...string) (string, string, error, string) {
	return RunStdin(workDir, "", args...) 
}

func RunStdin(workDir, stdin string, args ...string) (string, string, error, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var stdoutReturnString string = ""
	var timeoutReached bool = false

	// measure execution time
	start := time.Now()

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
	case <-time.After(30 * time.Second):
		err := cmd.Process.Kill()
		if err != nil {
			log.Println("2. failed to kill process\n", err.Error())
		} 
		timeoutReached = true
		log.Println("ERROR: process killed as timeout reached")
	case err := <-done:
		if err != nil {
			// default error message
			errMsg := err.Error()
			if exitErr, ok := err.(*exec.ExitError); ok {
				// only works on Unix systems
				// added this code to get `segmentation fault` for C++ to return
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					if status.Signaled() {
						sig := status.Signal()
						errMsg = "terminated by signal: " + sig.String()
						log.Printf("3. process finished due to signal: %s\n", sig.String())
					} else {
						log.Printf("3. process exited with code: %d\n", status.ExitStatus())
					}
				}
			}
			err = errors.New(errMsg)
			return stdout.String(), stderr.String(), err, "1"
		}
	}

	// format example: 186
	duration := time.Since(start)
	milliseconds := duration.Milliseconds()
	durationString := strconv.FormatInt(milliseconds, 10)

	// if timeout occured, prevent stdout from returning
	if timeoutReached {
		stdoutReturnString = ""
		err = errors.New("process killed as timeout reached")
	} else {
		// do not return a stdout string longer than N characters
		stdoutReturnString = stdout.String()
		if (len(stdoutReturnString) > 20000) {
			stdoutReturnString = stdoutReturnString[0:20000]
		} 
	}

	return stdoutReturnString, stderr.String(), err, durationString
}

func RunBash(workDir, command string) (string, string, error, string) {
	return Run(workDir, "bash", "-c", command)
}

func RunBashStdin(workDir, command, stdin string) (string, string, error, string) {
	return RunStdin(workDir, stdin, "bash", "-c", command)
}
