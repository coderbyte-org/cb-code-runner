package language

import (
	"github.com/danielborowski/cb-code-runner/language/javascript"
	"github.com/danielborowski/cb-code-runner/language/python"
	"github.com/danielborowski/cb-code-runner/language/ruby"
	"github.com/danielborowski/cb-code-runner/language/php"
	"github.com/danielborowski/cb-code-runner/language/java"
	"github.com/danielborowski/cb-code-runner/language/swift"
	"github.com/danielborowski/cb-code-runner/language/golang"
	"github.com/danielborowski/cb-code-runner/language/cpp"
	"github.com/danielborowski/cb-code-runner/language/csharp"
	"github.com/danielborowski/cb-code-runner/language/c"
	"github.com/danielborowski/cb-code-runner/language/kotlin"
)

type runFn func([]string, string) (string, string, error)

var languages = map[string]runFn{
	"javascript":   javascript.Run,
	"python":       python.Run,
	"ruby":         ruby.Run,
	"php":          php.Run,
	"java":         java.Run,
	"swift":        swift.Run,
	"go":           golang.Run,
	"cpp":          cpp.Run,
	"csharp":       csharp.Run,
	"c":            c.Run,
	"kotlin":       kotlin.Run,
}

func IsSupported(lang string) bool {
	_, supported := languages[lang]
	return supported
}

func Run(lang string, files []string, stdin string) (string, string, error) {
	return languages[lang](files, stdin)
}
