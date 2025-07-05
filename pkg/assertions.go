package specform

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/specform/go-sdk/types"
)

type AssertionFn func(value string, output string, ctx *types.AssertionContext) types.AssertionResult

type AssertionRegistry struct {
	registry map[string]AssertionFn
}

/**
 * Factory function to create a new AssertionRegistry
 */
func NewAssertionRegistry() *AssertionRegistry {
	return &AssertionRegistry{
		registry: make(map[string]AssertionFn),
	}
}

func (r *AssertionRegistry) Register(name string, fn AssertionFn) error {
	if _, exists := r.registry[name]; exists {
		return fmt.Errorf("Assertion %s already registered", name)
	}
	r.registry[name] = fn
	return nil
}

func (r *AssertionRegistry) Get(name string) (AssertionFn, error) {
	if fn, exists := r.registry[name]; exists {
		return fn, nil
	}
	return nil, fmt.Errorf("Assertion %s not found", name)
}

func (r *AssertionRegistry) Has(name string) bool {
	if _, exists := r.registry[name]; exists {
		return true
	}
	return false
}

func (r *AssertionRegistry) Run(name, value, output string, ctx *types.AssertionContext) (types.AssertionResult, error) {
	fn, err := r.Get(name)
	if err != nil {
		return types.AssertionResult{
			Type:    name,
			Value:   value,
			Passed:  false,
			Message: "✘ " + err.Error(),
		}, err
	}
	return fn(value, output, ctx), nil
}

func (r *AssertionRegistry) RunAll(output string, assertions []types.Assertion, ctx *types.AssertionContext) []types.AssertionResult {
	results := make([]types.AssertionResult, 0, len(assertions))

	for _, a := range assertions {
		result, err := r.Run(a.Type, a.Value, output, ctx)
		if err != nil {
			results = append(results, types.AssertionResult{
				Type:    a.Type,
				Value:   a.Value,
				Passed:  false,
				Message: "✘ " + err.Error(),
			})
			continue
		}
		results = append(results, result)
	}

	return results
}

var defaultRegistry = initDefaultRegistry()

func initDefaultRegistry() *AssertionRegistry {
	r := NewAssertionRegistry()

	// contains
	r.Register("contains", func(value, output string, _ *types.AssertionContext) types.AssertionResult {
		// Normalize the output and value
		normalizedOutput := normalizeText(output)
		normalizedValue := normalizeText(value)

		// Check if the normalized output contains the normalized value
		passed := strings.Contains(normalizedOutput, normalizedValue)
		msg := passFailMsg(passed, "Output contains '%s'", "Output missing '%s'", value)
		return types.AssertionResult{Type: "contains", Value: value, Passed: passed, Message: msg}
	})

	// equals
	r.Register("equals", func(value, output string, _ *types.AssertionContext) types.AssertionResult {
		// Normalize the output and value
		normalizedOutput := normalizeText(output)
		normalizedValue := normalizeText(value)
		// Check if the normalized output equals the normalized value
		passed := strings.TrimSpace(normalizedOutput) == strings.TrimSpace(normalizedValue)
		msg := passFailMsg(passed, "Output exactly matches expected value", "Output does not match expected value", "")
		return types.AssertionResult{Type: "equals", Value: value, Passed: passed, Message: msg}
	})

	// matches
	r.Register("matches", func(value, output string, _ *types.AssertionContext) types.AssertionResult {

		pattern, flags := parseRegex(value)

		var re *regexp.Regexp
		var err error
		if strings.Contains(flags, "i") {
			re, err = regexp.Compile("(?i)" + pattern)
		} else {
			re, err = regexp.Compile(pattern)
		}

		if err != nil {
			return types.AssertionResult{Type: "matches", Value: value, Passed: false, Message: fmt.Sprintf("✘ Invalid regex: %s", err)}
		}

		passed := re.MatchString(output)
		msg := passFailMsg(passed, "Output matches regex %s", "Output does not match regex %s", value)
		return types.AssertionResult{Type: "matches", Value: value, Passed: passed, Message: msg}
	})

	// semantic-similarity
	r.Register("semantic-similarity", func(value, _ string, ctx *types.AssertionContext) types.AssertionResult {
		score := 0.0
		if ctx != nil && ctx.SemanticScores != nil {
			if s, ok := ctx.SemanticScores[value]; ok {
				score = s
			}
		}
		threshold := 0.85
		if ctx != nil && ctx.Threshold > 0 {
			threshold = ctx.Threshold
		}
		passed := score >= threshold
		msg := fmt.Sprintf("%s semantic similarity %.2f vs threshold %.2f", boolPrefix(passed), score, threshold)
		return types.AssertionResult{Type: "semantic-similarity", Value: value, Passed: passed, Message: msg}
	})

	return r
}

// RunAssertions executes a list of assertions against the provided output string.
//
// Parameters:
//   - output: The string output to validate against the assertions.
//   - assertions: A slice of Assertion objects defining the checks to perform.
//   - ctx: An optional AssertionContext providing additional context for the assertions.
//
// Returns:
//
//	A slice of AssertionResult objects, each representing the result of an assertion.
func RunAssertions(output string, assertions []types.Assertion, ctx *types.AssertionContext) []types.AssertionResult {
	return defaultRegistry.RunAll(output, assertions, ctx)
}

/**
 * Public API to run a single assertion
 */
func RunAssertion(output string, assertion types.Assertion, ctx *types.AssertionContext) (types.AssertionResult, error) {
	return defaultRegistry.Run(assertion.Type, assertion.Value, output, ctx)
}

/**
 * Public API to register a new assertion
 */
func RegisterAssertion(name string, fn AssertionFn) error {
	return defaultRegistry.Register(name, fn)
}

func normalizeText(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)
	return s
}

func parseRegex(pattern string) (string, string) {
	if strings.HasPrefix(pattern, "/") && strings.LastIndex(pattern, "/") > 0 {
		lastSlash := strings.LastIndex(pattern, "/")
		return pattern[1:lastSlash], pattern[lastSlash+1:]
	}
	return pattern, ""
}

func passFailMsg(passed bool, okFmt, failFmt, val string) string {
	if passed {
		return "✔ " + fmt.Sprintf(okFmt, val)
	}
	return "✘ " + fmt.Sprintf(failFmt, val)
}

func boolPrefix(ok bool) string {
	if ok {
		return "✔"
	}
	return "✘"
}
