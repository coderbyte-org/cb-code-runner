cb-code-runner
================
Clone of [glot.io](https://github.com/prasmussen/glot). Modified to combine server + code runner into one application.

## Overview
Application that reads json payload from request body, executes the code, and return json output.

## Testing/Building
To test the Go code execution program, first install dependencies with `go get -d ./...` then start the server:

```
PORT=8085 go run runner.go
```

To build the Go binary, run:

```
bash build_linux.sh
```

## Examples
The input JSON payload expects two parameters: `language` and a `files` array with each object in `files` containing a file `name` and code `content`.

##### Input
```javascript
{
  "language": "javascript",
  "files": [
    {
      "name": "main.js",
      "content": "console.log(2+2+2);"
    }
  ]
}
```

##### Output
```javascript
{
  "stdout": "6\n",
  "stderr": "",
  "error": ""
}
```