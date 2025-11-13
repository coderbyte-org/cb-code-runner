package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type CbConfig struct {
	Run string `json:"run"`
}

// ParseCbConfig locates ".cbconfig", reads it, parses JSON, and
// returns a run command slice such as ["php", "main.php"].
// Returns nil if missing or invalid.
func ParseCbConfig(files []string) []string {
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

	run := strings.TrimSpace(cfg.Run)
	if run == "" {
		return nil
	}

	parts := strings.Fields(run)
	if len(parts) == 0 {
		return nil
	}

	return parts
}