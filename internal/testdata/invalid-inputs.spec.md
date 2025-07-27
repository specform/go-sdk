---
model: "gpt-4"
scenario: "Summarize a technical article"
temperature: 0.3
tags: ["summarization", "test"]
---

```prompt
Please summarize this article: {{article}}
Make sure to use a {{tone}} tone.
```

```inputs
this line has no equals
this one isn't quoted =
another broken one =
```

```assertions
- contains: "real time"
- matches: /HTTP/i
- semantic-similarity: "event-driven communication"
```
