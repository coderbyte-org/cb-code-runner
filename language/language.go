package language

import (
	"github.com/danielborowski/cb-code-runner/language/javascript"
)

type runFn func([]string, string) (string, string, error)

var languages = map[string]runFn{
	"javascript":   javascript.Run,
}

func IsSupported(lang string) bool {
	_, supported := languages[lang]
	return supported
}

func Run(lang string, files []string, stdin string) (string, string, error) {
	return languages[lang](files, stdin)
}