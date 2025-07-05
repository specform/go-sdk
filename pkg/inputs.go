package specform

import (
	"bufio"
	"strings"
)

func ParseInputBlock(content string) ([]string, map[string]string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	vars := []string{}
	defaults := make(map[string]string)

	var currentKey string
	var currentVal strings.Builder
	inMultiline := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// This is a bit of a nieve way to do this, but it works
		// for our purposes. We can always improve this later with a
		// more robust parser if needed.
		switch {
		// Multiline string start
		case strings.HasSuffix(line, "=\"\"\""):
			inMultiline = true
			// Grab our key from the left side of the equals sign
			currentKey = strings.Split(line, "=")[0]
			// Reset and move on to the next line
			currentVal.Reset()

		// Multiline string end
		case line == "\"\"\"" && inMultiline:
			inMultiline = false
			// Add the current key and value to the map
			vars = append(vars, currentKey)

		// We're in a multiline string, so we need to
		// append the line to the current value
		case inMultiline:
			currentVal.WriteString(line + "\n")

		// Normal line, so we need to check if it contains an equals sign
		// and split on the first equals sign. If it does, this input variable
		// has a default value set
		case strings.Contains(line, "="):
			parts := strings.SplitN(line, "=", 2)
			// Set the current key and value
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			vars = append(vars, key)

			// Set the default value of the input variable
			defaults[key] = strings.Trim(val, "\"")

		// If the line is not empty and does not contain an equals sign,
		// we need to add it to the list of variables without a default value
		case line != "":
			vars = append(vars, line)
		}
	}

	return vars, defaults
}
