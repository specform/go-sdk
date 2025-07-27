---
model: "gpt-4"
title: "Summarize a technical article"
temperature: 0.3
tags: ["summarization", "test"]
---

## Description

This is a test `promptspec` for testing our parser. This description should be ignored and not returned in the final object.

## Prompt

This text should be ignored in the final object. This allows us to document the prompt without it being include in the final object

```prompt
Please summarize this article: {{article}}
Make sure to use a {{tone}} tone.
```

```inputs
article = """Webhooks enable real-time communication..."""
tone = "casual"
```

```assertions
- contains: "real time"
- matches: /HTTP/i
- semantic-similarity: "event-driven communication"
```
