---
model: "gpt-4"
scenario: "Summarize a technical article"
temperature: 0.3
tags: ["summarization", "test"]
---

```inputs
article = """Webhooks enable real-time communication..."""
tone = "casual"
```

```assertions
- contains: "real time"
- matches: /HTTP/i
- semantic-similarity: "event-driven communication"
```
