package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents the key-value pairs parsed from a .env file.
type EnvMap map[string]string

// ParseFile reads a .env file from the given path and returns an EnvMap.
// It skips blank lines and comment lines (starting with '#').
// It returns an error if the file cannot be opened or a line is malformed.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("parser: malformed line %d in %q: missing '='", lineNum, path)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Strip optional surrounding quotes from value
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: error reading %q: %w", path, err)
	}

	return env, nil
}
