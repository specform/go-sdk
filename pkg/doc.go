// Package specform provides tools for compiling, rendering, and asserting prompts.
//
// This SDK is the Go implementation of the specform prompt framework. It allows you to:
//
//   - Compile .spec.md files into structured prompt definitions
//   - Render compiled prompts with dynamic inputs
//   - Register and evaluate assertions on model output
//   - Work with snapshots for regression and evaluation
//
// Common usage:
//
//	// Compile multiple .spec.md files
//	prompt, err := specform.CompileSpecFiles("path/to/files", "output/dir")
//
//	// Render a prompt with inputs
//	output, err := specform.RenderPrompt(prompt, map[string]string{"input": "value"})
//
//	// Run assertions on model output
//	results := specform.RunAssertions(output, prompt.Assertions, nil)
//
//	// Register a custom assertion type
//	specform.RegisterAssertion("starts-with", func(val, out string, ctx *types.AssertionContext) types.AssertionResult { ... })
//
//	// Save a snapshot of the model output
//	err := specform.SaveSnapshot("path/to/snapshot.json", prompt, output, results, inputs)
//	// Load a snapshot for comparison
//	snapshot, err := specform.LoadSnapshot("path/to/snapshot.json")
//
// // CosineSimilarity is a common metric for comparing text similarity. The SDK provides
// // a function to calculate cosine similarity scores between two text inputs.
// // This can be used in assertions to check if the model output is semantically similar
// // to a reference output.
// specform.CosineSimilarity("text1", "text2")
package specform
