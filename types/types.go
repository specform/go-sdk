package types

import "time"

type Assertion struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CompiledPrompt struct {
	ID          string            `json:"id"`
	Hash        string            `json:"hash"`
	Feature     string            `json:"feature,omitempty"`
	Scenario    string            `json:"scenario"`
	Prompt      string            `json:"compiledPrompt"`
	Inputs      []string          `json:"inputs"`
	Values      map[string]string `json:"defaultInputs"`
	Assertions  []Assertion       `json:"assertions,omitempty"`
	Snapshot    string            `json:"snapshot,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Model       string            `json:"model"`
	Temperature float64           `json:"temperature,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	SourcePath  string            `json:"sourcePath,omitempty"`
}

type Snapshot struct {
	ID         string            `json:"id"`
	Hash       string            `json:"hash"`
	Output     string            `json:"output"`
	Inputs     map[string]string `json:"inputs"`
	Assertions []AssertionResult `json:"assertions"`
	Passed     bool              `json:"passed"`
	Timestamp  time.Time         `json:"timestamp"`
}

type AssertionResult struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

// AssertionContext holds optional external data used to evaluate advanced assertions.
type AssertionContext struct {
	SemanticScores map[string]float64 // expected value → similarity score
	Threshold      float64            // override threshold (default 0.85)
}
