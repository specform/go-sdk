package specform

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/specform/go-sdk/types"
)

type RenderOptions struct {
	Strict bool // If true, all variables are to be set
}

func RenderPrompt(prompt *types.CompiledPrompt, inputs map[string]string, opts *RenderOptions) (string, error) {
	// Normalize {{var}} → {{.var}} for Go template engine
	varPattern := regexp.MustCompile(`{{\s*(\w+)\s*}}`)
	normalizedPrompt := varPattern.ReplaceAllString(prompt.Prompt, "{{.$1}}")

	tpl, err := template.New("prompt").Option("missingkey=error").Parse(normalizedPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	// Set default inputs from the scenario if they are not overridden
	// by the user
	merged := map[string]string{}
	for k, v := range prompt.Values {
		merged[k] = v
	}
	// Merge user inputs into the merged map
	// This will override any default values
	for k, v := range inputs {
		merged[k] = v
	}

	// Validate that all require inputs are set (strict mode)
	if opts != nil && opts.Strict {
		missing := []string{}
		for _, input := range prompt.Inputs {
			if _, ok := merged[input]; !ok {
				missing = append(missing, input)
			}
		}

		if len(missing) > 0 {
			return "", fmt.Errorf("missing required inputs: %s", strings.Join(missing, ", "))
		}
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, merged); err != nil {
		return "", fmt.Errorf("failed to renderprompt: %w", err)
	}

	return buf.String(), nil
}
