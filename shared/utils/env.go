package utils

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	LoadEnv()
}

// LoadEnv searches for .env starting from current working directory upwards (up to 5 levels)
// and loads its values into the environment. If it fails, it logs the problem and ensures the app continues.
func LoadEnv() {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("[WARNING] LoadEnv: could not determine working directory: %v\n", err)
		return
	}

	var envFile string
	for i := 0; i < 5; i++ {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			envFile = path
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if envFile == "" {
		log.Println("[WARNING] LoadEnv: .env file not found in current or parent directories. Falling back to default values.")
		return
	}

	file, err := os.Open(envFile)
	if err != nil {
		log.Printf("[WARNING] LoadEnv: failed to open .env file: %v. Falling back to default values.\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("[WARNING] LoadEnv: invalid line format at %s:%d: %q. Skipping.\n", envFile, lineNum, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Strip quotes if present
		if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
			(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
			if len(val) >= 2 {
				val = val[1 : len(val)-1]
			}
		}

		// Set environment variable if not already set. If already set, do not overwrite.
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, val); err != nil {
				log.Printf("[WARNING] LoadEnv: failed to set env var %s: %v\n", key, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[WARNING] LoadEnv: error reading .env file: %v. Falling back to default values.\n", err)
	} else {
		log.Printf("[INFO] LoadEnv: environment loaded successfully from %s\n", envFile)
	}
}

// GetEnv retrieves an environment variable or returns a fallback value if not set.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

