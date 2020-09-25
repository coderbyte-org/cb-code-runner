package language

import (
	"github.com/coderbyte-org/cb-code-runner/language/javascript"
	"github.com/coderbyte-org/cb-code-runner/language/python"
	"github.com/coderbyte-org/cb-code-runner/language/ruby"
	"github.com/coderbyte-org/cb-code-runner/language/php"
	"github.com/coderbyte-org/cb-code-runner/language/java"
	"github.com/coderbyte-org/cb-code-runner/language/swift"
	"github.com/coderbyte-org/cb-code-runner/language/golang"
	"github.com/coderbyte-org/cb-code-runner/language/cpp"
	"github.com/coderbyte-org/cb-code-runner/language/csharp"
	"github.com/coderbyte-org/cb-code-runner/language/c"
	"github.com/coderbyte-org/cb-code-runner/language/kotlin"
	"github.com/coderbyte-org/cb-code-runner/language/typescript"
	"github.com/coderbyte-org/cb-code-runner/language/clojure"
	"github.com/coderbyte-org/cb-code-runner/language/bash"
	"github.com/coderbyte-org/cb-code-runner/language/elixir"
	"github.com/coderbyte-org/cb-code-runner/language/scala"
	"github.com/coderbyte-org/cb-code-runner/language/rust"
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
	"typescript":   typescript.Run,
	"clojure":      clojure.Run,
	"bash":         bash.Run,
	"elixir":       elixir.Run,
	"scala":        scala.Run,
	"rust":         rust.Run,
}

func IsSupported(lang string) bool {
	_, supported := languages[lang]
	return supported
}

func Run(lang string, files []string, stdin string) (string, string, error) {
	return languages[lang](files, stdin)
}
