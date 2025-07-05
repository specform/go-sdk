# Specform (Go SDK + CLI)

**`specform`** is a prompt specification framework that supports prompt compilation, rendering, assertions, and snapshot testing.

This package includes both:

- A Go SDK (`pkg/specform`) for embedding prompt logic in Go programs
- A CLI (`specform`) for working with prompt specs on the command line

---

## Features

- Compile `.spec.md` files into structured JSON prompts
- Render prompts with dynamic input values
- Define and evaluate assertions on model outputs
- Create and test snapshots for regression testing
- Serve compiled prompts and snapshots over HTTP
- Register custom assertion types (Go SDK only)

---

## 📦 Install (CLI)

```bash
go install github.com/specform/specform/specform@latest
```

Or build from source:

```bash
cd sdk/go
go build -o specform ./cmd/specform
```

---

## 🥪 CLI Usage

### Compile

```bash
specform compile ./examples --output build
```

Options:

- `--watch` – Watch files for changes
- `--stdout` – Output JSON to stdout instead of files

---

### Render

```bash
specform render --prompt build/my-prompt.prompt.json --input name=Alice
```

You can also pass `--inputs inputs.json`.

---

### Test

```bash
specform test --prompt build/my-prompt.prompt.json --output output.txt
```

Optional:

- `--similarity scores.json` – Provide semantic similarity scores
- `--inputs` / `--input` for variable values

---

### Snapshot

```bash
specform snapshot --prompt build/my-prompt.prompt.json --output output.txt --out snapshots/
```

Only saves the snapshot if all assertions pass.

---

### Serve

```bash
specform serve --dir build --port 8080
```

Serves:

- `/prompts` and `/prompts/:id`
- `/snapshots` and `/snapshots/:id`

---

## Go SDK Usage

### Compile

```go
prompt, err := specform.CompileSpecFiles("examples/hello.spec.md", nil)
```

### Render

```go
output, err := specform.RenderPrompt(prompt, map[string]string{"name": "Alice"})
```

### Assert

```go
results := specform.RunAssertions(output, prompt.Assertions, nil)
```

### Register custom assertion

```go
specform.RegisterAssertion("starts-with", func(val, out string, _ *types.AssertionContext) types.AssertionResult {
  passed := strings.HasPrefix(out, val)
  return types.AssertionResult{Type: "starts-with", Value: val, Passed: passed}
})
```

---

### Snapshot Helpers

```go
err := specform.SaveSnapshot("snapshots/hello.snap.json", prompt, output, results, inputs)

snap, err := specform.LoadSnapshot("snapshots/hello.snap.json")
```

---

## Project Structure

```
sdk/go/
├── pkg/          # Go SDK (public)
├── cmd/specform  # CLI entry point
├── internal/     # Internal logic (parser, compiler)
├── types/        # Shared types
```

---

## Design Notes

- Go SDK is designed to be **minimal and idiomatic**
- TypeScript SDK uses `createClient()` for stateful use
- Go SDK uses **package-level functions** and optional extensibility (e.g. assertion registration)

---

## Tests

Run:

```bash
go test ./...
```

---

## 📄 License

© 2025 Brad Walker. Licensed under the Apache License, Version 2.0, see [LICENSE](./LICENSE) for details.
