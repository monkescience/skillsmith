---
name: code-principles
description: Code design principles covering style, naming, error handling, and testing patterns. Load when writing or reviewing code to follow established conventions.
compatibility:
  - opencode
  - claude
license: MIT
---

## Design

- Prefer library-provided utilities over custom implementations
- Prefer small, focused functions over large monolithic ones
- Prefer explicit over implicit behavior
- Prefer composition over inheritance
- Prefer immutability where practical
- Prefer simple code over unnecessary abstractions
- Prefer changing code directly over adding backwards-compat layers

## Style

- Run formatters before committing; do not manually format code
- Follow language-specific conventions
- Do not change formatting in code you are not modifying

## Naming & Comments

- Prefer descriptive names over comments
- Avoid comments unless they add real value
- Prefer deleting dead code over commenting it out

## Task Tracking with TODOs

- When discovering missing or incomplete implementations, insert TODO comments directly in the relevant code
- TODOs must be explicit and scoped: `// TODO: implement retry logic with exponential backoff`
- Never leave vague TODOs like `// TODO: fix this` or `// TODO: implement`
- When deferring work, explain why and what's needed
- Prefer embedding intent in code over external documentation
- Remove TODOs only when the work is complete

## Error Handling

- Handle errors at appropriate boundaries (API, user input, external calls)
- Propagate errors with context rather than swallowing them
- Prefer returning errors over panicking/throwing
- Fail fast on programmer errors; handle gracefully on user/external errors
- Do not add defensive checks for conditions the type system already prevents

## Testing

- Use given-when-then pattern:
  - given: setup and preconditions
  - when: action being tested
  - then: expected outcome
- Use comments to mark each section when helpful
