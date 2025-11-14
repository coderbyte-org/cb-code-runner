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
	"io"
	"os"
)

// helper function because we treat TS errors slightly differently
func isTypescriptRun(args []string) bool {
	if len(args) == 0 {
		return false
	}
	cmd := args[0]
	// Heuristics: mark as TS if we’re using ts-node, tsc, Jest/ts-jest, etc.
	if strings.Contains(cmd, "ts-node") || strings.Contains(cmd, "tsc") || strings.Contains(cmd, "jest") {
		return true
	}
	// Also treat commands that run .ts/.tsx files directly as TS
	for _, a := range args[1:] {
		if strings.HasSuffix(a, ".ts") || strings.HasSuffix(a, ".tsx") {
			return true
		}
	}
	return false
}

func Run(workDir string, args ...string) (string, string, error, string) {
	return RunStdin(workDir, "", args...) 
}

// Sentinel used to stop io.Copy when stdout cap is reached
var errStdoutLimitReached = errors.New("stdout limit reached")

// limitedWriter writes up to 'limit' bytes, then returns errStdoutLimitReached.
// It still stores what it wrote so far in the internal buffer.
type limitedWriter struct {
	buf   *bytes.Buffer
	limit int
	wrote int
}

func (lw *limitedWriter) Write(p []byte) (int, error) {
	remain := lw.limit - lw.wrote
	if remain <= 0 {
		return 0, errStdoutLimitReached
	}
	if len(p) > remain {
		// write only the allowed prefix, then signal limit
		n, _ := lw.buf.Write(p[:remain])
		lw.wrote += n
		return n, errStdoutLimitReached
	}
	n, _ := lw.buf.Write(p)
	lw.wrote += n
	return n, nil
}

func RunStdin(workDir, stdin string, args ...string) (string, string, error, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var stdoutReturnString string = ""
	var timeoutReached bool = false
	var err error

	const maxStdoutBytes = 30000 // show first N kb then kill
	const wallTimeout = 30 * time.Second

	start := time.Now()

	// If running Python, force unbuffered so prints appear immediately.
	// You can also pass "-u" in args at the caller; this env var is a safe default.
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "PYTHONUNBUFFERED=1")
	cmd.Stdin = strings.NewReader(stdin)

	stdoutPipe, pipeErr1 := cmd.StdoutPipe()
	if pipeErr1 != nil {
		log.Println("stdout pipe error:", pipeErr1)
	}
	stderrPipe, pipeErr2 := cmd.StderrPipe()
	if pipeErr2 != nil {
		log.Println("stderr pipe error:", pipeErr2)
	}

	// Start process AFTER creating pipes
	err = cmd.Start()
	if err != nil {
		log.Println("1. process fail\n", err.Error())
	}

	// Detect if this looks like a TypeScript-related command
	isTS := isTypescriptRun(args)

	// We'll copy stdout into a limited writer (cap), and stderr into a normal buffer.
	// When the cap is hit, we kill the process immediately.
	limitHitCh := make(chan struct{}, 1)

	lw := &limitedWriter{buf: &stdout, limit: maxStdoutBytes}

	go func() {
		_, copyErr := io.Copy(lw, stdoutPipe)
		if copyErr == errStdoutLimitReached {
			// Cap reached → signal and kill process
			select { case limitHitCh <- struct{}{}: default: }
		}
	}()

	go func() {
		_, _ = io.Copy(&stderr, stderrPipe)
	}()

	// Waiter
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	// Main loop: stop on (1) cap reached, (2) process exit, or (3) wall timeout
	timer := time.NewTimer(wallTimeout)
	defer timer.Stop()

	killedForOutputLimit := false

	loop:
		for {
			select {
			case <-limitHitCh:
				// Kill fast once stdout cap is reached
				_ = cmd.Process.Kill()
				killedForOutputLimit = true
				break loop
			case werr := <-done:
				// Process exited on its own
				if werr != nil {
					if isTS {
						// TS PATH: log but DO NOT propagate error or return early
						log.Printf("TS process exited with error (suppressed): %v\n", werr)
					} else {
						// Build a better error message (your original logic)
						errMsg := werr.Error()
						if exitErr, ok := werr.(*exec.ExitError); ok {
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
				// Normal exit
				break loop
			case <-timer.C:
				_ = cmd.Process.Kill()
				timeoutReached = true
				break loop
			}
		}

	// Duration string
	duration := time.Since(start)
	milliseconds := duration.Milliseconds()
	durationString := strconv.FormatInt(milliseconds, 10)

	// Post-process outcomes
	if killedForOutputLimit {
		// Append a helpful trailer for the user.
		// Note: we add the trailer on top of the captured stdout prefix.
		stdoutReturnString = stdout.String()
		trailer := "\n...\n(process killed due to output limit; showing first " +
			strconv.Itoa(maxStdoutBytes) + " bytes of stdout)"
		stdoutReturnString += trailer

		// Keep stderr as-is, return a specific error
		err = errors.New("E_OUTPUT_CAP: process killed after reaching stdout limit")
		// You can return a more specific code in the 4th position if desired; keeping "1" for compatibility.
		return stdoutReturnString, stderr.String(), err, "1"
	}

	if timeoutReached {
		// You currently blank stdout when timeout occurs; keep your behavior or
		// switch to returning the partial buffer if you prefer.
		err = errors.New("process killed as timeout reached")
		return "", stderr.String(), err, durationString
	}

	// Normal completion path: enforce your existing 20k cap
	stdoutReturnString = stdout.String()
	if len(stdoutReturnString) > 20000 {
		stdoutReturnString = stdoutReturnString[0:20000]
	}
	return stdoutReturnString, stderr.String(), err, durationString
}

func RunBash(workDir, command string) (string, string, error, string) {
	return Run(workDir, "bash", "-c", command)
}

func RunBashStdin(workDir, command, stdin string) (string, string, error, string) {
	return RunStdin(workDir, stdin, "bash", "-c", command)
}
