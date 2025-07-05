package internal

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/specform/go-sdk/types"
)

func ParseAssertionsBlock(content string) ([]types.Assertion, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var out []types.Assertion

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "-") {
			continue
		}

		line = strings.TrimPrefix(line, "-")
		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			continue
		}

		typName := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		out = append(out, types.Assertion{
			Type:  typName,
			Value: strings.Trim(val, "\""),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse assertions block: %w", err)
	}

	return out, nil
}
