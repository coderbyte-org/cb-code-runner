package cpp

import (
	"github.com/coderbyte-org/cb-code-runner/config"
	"github.com/coderbyte-org/cb-code-runner/cmd"
	"github.com/coderbyte-org/cb-code-runner/util"
	"path/filepath"
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func resolveAddr2line(binPath, addr string) (string, error) {
	cmd := exec.Command("addr2line", "-f", "-e", binPath, addr)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func EnhanceStackTrace(stderr, binPath string) string {
	if stderr == "" || binPath == "" {
		return stderr
	}

	var enhanced strings.Builder
	enhanced.WriteString(stderr)

	scanner := bufio.NewScanner(strings.NewReader(stderr))
	re := regexp.MustCompile(`\[(0x[0-9a-fA-F]+)\]`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			addr := matches[1]
			loc, err := resolveAddr2line(binPath, addr)
			if err == nil {
				enhanced.WriteString(fmt.Sprintf("  ↳ %s\n", loc))
			}
		}
	}

	return enhanced.String()
}

func Run(files []string, stdin string) (string, string, error, string) {
	workDir := filepath.Dir(files[0])
	binName := "a.out"

	sourceFiles := util.FilterByExtension(files, "cpp")

	// Default compile command:
	compileArgs := append([]string{"clang++", "-std=c++20", "-g", "-no-pie", "-fno-pie", "-O0", "-o", binName}, sourceFiles...)

	// If .cbconfig has "compile", use it instead, with support for ENV_SOURCE_FILES
	if cfgCompile := config.ParseCbConfigField(files, "compile"); cfgCompile != nil {
		var expanded []string
		for _, part := range cfgCompile {
			if part == "ENV_SOURCE_FILES" {
				expanded = append(expanded, sourceFiles...)
			} else {
				expanded = append(expanded, part)
			}
		}
		// If there is no ENV_SOURCE_FILES in the config, we assume the user
		// fully specified the sources themselves and we do NOT auto-append.
		compileArgs = expanded
	}

	stdout, stderr, err, duration := cmd.Run(workDir, compileArgs...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	// compilation succeeded, now run binary
	binPath := filepath.Join(workDir, binName)

	// Default run: execute the compiled binary directly
	runArgs := []string{binPath}

	// Allow override via .cbconfig "run"
	// e.g. "run": "./a.out" or "run": "valgrind ./a.out"
	if cfgRun := config.ParseCbConfigField(files, "run"); cfgRun != nil {
		runArgs = cfgRun
	}

	stdout, stderr, err, duration = cmd.RunStdin(workDir, stdin, runArgs...)
	
	if err != nil || stderr != "" {
		stderr = EnhanceStackTrace(stderr, binPath)
	}

	return stdout, stderr, err, duration
}