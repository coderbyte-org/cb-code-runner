package cpp

import (
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
	args := append([]string{"clang++", "-std=c++20", "-g", "-no-pie", "-fno-pie", "-O0", "-o", binName}, sourceFiles...)
	stdout, stderr, err, duration := cmd.Run(workDir, args...)
	if err != nil || stderr != "" {
		return stdout, stderr, err, duration
	}

	// compilation succeeded, now run binary
	binPath := filepath.Join(workDir, binName)
	stdout, stderr, err, duration = cmd.RunStdin(workDir, stdin, binPath)
	
	if err != nil || stderr != "" {
		stderr = EnhanceStackTrace(stderr, binPath)
	}

	return stdout, stderr, err, duration
}