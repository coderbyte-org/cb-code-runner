package language

import (
	"github.com/danielborowski/cb-code-runner/language/javascript"
	"github.com/danielborowski/cb-code-runner/language/python"
	"github.com/danielborowski/cb-code-runner/language/ruby"
)

type runFn func([]string, string) (string, string, error)

var languages = map[string]runFn{
	"javascript":   javascript.Run,
	"python":       python.Run,
	"ruby":         ruby.Run,
}

func IsSupported(lang string) bool {
	_, supported := languages[lang]
	return supported
}

func Run(lang string, files []string, stdin string) (string, string, error) {
	return languages[lang](files, stdin)
}
