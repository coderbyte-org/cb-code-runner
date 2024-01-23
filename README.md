cb-code-runner
================
Application that reads json payload from request body, executes the code, and return json output. Clone of [glot.io](https://github.com/prasmussen/glot). Modified to combine server + code runner into one application.

## Testing
To test the Go code execution program, first install dependencies with `go get` then start the server:

```
PORT=8085 go run runner.go
```

Then test the program by running:

```
curl -X POST -d "{\"language\": \"javascript\",\"files\": [{\"name\": \"main.js\",\"content\": \"console.log(2+2+2);\"}]}" http://localhost:8085
```

To build the Go binary, run:

```
bash build.sh
```

## Examples
The input JSON payload expects two parameters: `language` and a `files` array with each object in `files` containing a file `name` and code `content`.

The output will be a JSON payload with the following parameters: `stdout`, `stderr`, `error`, and `duration` which is the execution time in milliseconds.

##### Input
```javascript
{
  "language": "javascript",
  "files": [
    {
      "name": "main.js",
      "content": "console.log(2+2);"
    }
  ]
}
```

##### Output
```javascript
{
  "stdout": "4\n",
  "stderr": "",
  "error": "",
  "duration": "37"
}
```

## Languages Supported
* Apex
* Bash
* C
* C++
* C#
* Clojure
* Dart
* Elixir
* Go
* Java
* JavaScript
* Kotlin
* PHP
* Python
* R
* Ruby
* Rust
* Scala
* Swift
* TypeScript