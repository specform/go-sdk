package internal

import (
	"bufio"
	"fmt"
	"strings"
)

// ParseInputBlock parses the inputs block from a spec file and returns
// the list of input variables and their default values.
func ParseInputBlock(content string) ([]string, map[string]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	vars := []string{}
	defaults := make(map[string]string)

	var currentKey string
	var currentVal strings.Builder
	inMultiline := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		// Multiline string start
		case !inMultiline && strings.Contains(line, "=") && strings.Contains(line, "\"\"\""):
			// Parse the key
			parts := strings.SplitN(line, "=", 2)
			currentKey = strings.TrimSpace(parts[0])
			vars = append(vars, currentKey)
			currentVal.Reset()
			inMultiline = true

			// Process content after opening quotes, if any
			if len(parts) > 1 {
				contentPart := strings.TrimSpace(parts[1])
				if strings.HasPrefix(contentPart, "\"\"\"") {
					contentAfterQuotes := strings.TrimPrefix(contentPart, "\"\"\"")

					// Handle single-line triple-quoted string
					if strings.HasSuffix(contentAfterQuotes, "\"\"\"") {
						content := strings.TrimSuffix(contentAfterQuotes, "\"\"\"")
						defaults[currentKey] = content
						inMultiline = false
					} else {
						// First line of multiline content
						currentVal.WriteString(contentAfterQuotes + "\n")
					}
				}
			}

		// Standalone closing quotes
		case inMultiline && line == "\"\"\"":
			defaults[currentKey] = strings.TrimSuffix(currentVal.String(), "\n")
			inMultiline = false

		// Closing quotes with content on the same line
		case inMultiline && strings.HasSuffix(line, "\"\"\""):
			content := strings.TrimSuffix(line, "\"\"\"")
			currentVal.WriteString(content)
			defaults[currentKey] = strings.TrimSuffix(currentVal.String(), "\n")
			inMultiline = false

		// Inside a multiline string
		case inMultiline:
			currentVal.WriteString(line + "\n")

		// Single-line input with default value
		case strings.Contains(line, "="):
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			vars = append(vars, key)
			defaults[key] = strings.Trim(val, "\"")

		// Input without default value
		case line != "":
			key := strings.TrimSpace(line)
			vars = append(vars, key)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading input block: %w", err)
	}

	if inMultiline {
		return nil, nil, fmt.Errorf("unclosed multiline string for key: %s", currentKey)
	}

	return vars, defaults, nil
}
