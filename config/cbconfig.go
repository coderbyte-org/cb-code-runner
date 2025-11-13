package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type CbConfig struct {
	Run string `json:"run"`
	Compile string `json:"compile"`
}

// ParseCbConfigField looks for ".cbconfig" in files, and returns the
// chosen field ("run" or "compile") as a []string, e.g. ["node", "main.js"].
// Returns nil if the file doesn't exist, JSON is invalid, or the field is empty.
func ParseCbConfigField(files []string, field string) []string {
	var cfgPath string
	for _, f := range files {
		if filepath.Base(f) == ".cbconfig" {
			cfgPath = f
			break
		}
	}
	if cfgPath == "" {
		return nil
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil
	}

	var cfg CbConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}

	var raw string
	switch field {
	case "run":
		raw = cfg.Run
	case "compile":
		raw = cfg.Compile
	default:
		return nil
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return nil
	}

	return parts
}