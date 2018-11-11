package main

import (
	"encoding/json"
	"github.com/danielborowski/cb-code-runner/cmd"
	"github.com/danielborowski/cb-code-runner/language"
	"io/ioutil"
	"os"
	"path/filepath"
	"net/http"
	"log"
)

type Payload struct {
	Language string          `json:"language"`
	Files    []*InMemoryFile `json:"files"`
	Stdin    string          `json:"stdin"`
	Command  string          `json:"command"`
}

type InMemoryFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}

func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		var errOccured bool = false

		// get input data
		body, _ := ioutil.ReadAll(r.Body)
		bodyData := string(body)
		log.Println(bodyData)

		payload := &Payload{}
		err := json.Unmarshal([]byte(bodyData), payload)

		if err != nil {
			log.Println("Failed to parse input json (%s)\n", err.Error())
			errOccured = true
		}

		// Ensure that we have at least one file
		if !errOccured && len(payload.Files) == 0 {
			log.Println("No files given\n")
			errOccured = true
		}

		// Check if we support given language
		if !errOccured && !language.IsSupported(payload.Language) {
			log.Println("Language '%s' is not supported\n", payload.Language)
			errOccured = true
		}

		// Write files to disk
		filepaths, err := writeFiles(payload.Files)
		if !errOccured && err != nil {
			log.Println("Failed to write file to disk (%s)", err.Error())
			errOccured = true
		}

		var stdout, stderr string

		if (!errOccured) {
			if payload.Command == "" {
				stdout, stderr, err = language.Run(payload.Language, filepaths, payload.Stdin)
			} else {
				workDir := filepath.Dir(filepaths[0])
				stdout, stderr, err = cmd.RunBashStdin(workDir, payload.Command, payload.Stdin)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			writeResult(stdout, stderr, err, w)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Fatal(http.ListenAndServe(":" + port, nil))
}

// Writes files to disk, returns list of absolute filepaths
func writeFiles(files []*InMemoryFile) ([]string, error) {
	// Create temp dir
	tmpPath, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	paths := make([]string, len(files), len(files))
	for i, file := range files {
		path, err := writeFile(tmpPath, file)
		if err != nil {
			return nil, err
		}

		paths[i] = path

	}
	return paths, nil
}

// Writes a single file to disk
func writeFile(basePath string, file *InMemoryFile) (string, error) {
	// Get absolute path to file inside basePath
	absPath := filepath.Join(basePath, file.Name)

	// Create all parent dirs
	err := os.MkdirAll(filepath.Dir(absPath), 0775)
	if err != nil {
		return "", err
	}

	// Write file to disk
	err = ioutil.WriteFile(absPath, []byte(file.Content), 0664)
	if err != nil {
		return "", err
	}

	// Return absolute path to file
	return absPath, nil
}

func writeResult(stdout, stderr string, err error, writer http.ResponseWriter) {
	result := &Result{
		Stdout: stdout,
		Stderr: stderr,
		Error:  errToStr(err),
	}
	json.NewEncoder(os.Stdout).Encode(result)
	responseJSON, err := json.Marshal(result)
	writer.Write(responseJSON)
}

func errToStr(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}
