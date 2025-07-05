package specform

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/specform/go-sdk/types"
)

func SaveSnapshot(
	path string,
	scenario *types.CompiledPrompt,
	output string,
	results []types.AssertionResult,
	inputs map[string]string,
) error {
	snapshot := types.Snapshot{
		ID:         scenario.ID,
		Hash:       scenario.Hash,
		Output:     output,
		Inputs:     inputs,
		Assertions: results,
		Passed:     allAssertionsPassed(results),
		Timestamp:  time.Now(),
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create snapshot file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(snapshot); err != nil {
		return fmt.Errorf("failed to encode snapshot: %w", err)
	}

	return nil
}

func LoadSnapshot(path string) (*types.Snapshot, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open snapshot file: %w", err)
	}
	defer file.Close()

	var snapshot types.Snapshot
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&snapshot); err != nil {
		return nil, fmt.Errorf("failed to decode snapshot: %w", err)
	}
	return &snapshot, nil
}

func allAssertionsPassed(results []types.AssertionResult) bool {
	for _, result := range results {
		if !result.Passed {
			return false
		}
	}
	return true
}
