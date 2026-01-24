---
name: git-conventions
description: Git workflow conventions for Conventional Commits, atomic commits, and PR reviews with Conventional Comments. Load when working with git, commits, branches, or pull requests.
compatibility:
  - opencode
  - claude
license: MIT
---

## Commits

- Use Conventional Commits: `type(scope): description`
- Keep commits focused and atomic (one logical change per commit)
- Subject line only, no body
- Do not add Co-Authored-By lines

### Commit types

| Type | Triggers release |
|------|------------------|
| `feat`, `fix`, `perf` | Yes |
| `docs`, `style`, `refactor`, `test`, `chore`, `build`, `ci` | No |

### Breaking changes

Add `!` after type (or scope) to indicate a breaking change that affects external consumers.

- For services: API contract changes (endpoints, request/response formats, auth flows)
- For libraries: Public API changes (exported functions, types, interfaces)

When uncertain if a change is breaking, ask before committing.

### History rewriting

- Do not use `git commit --amend` on commits that have been pushed
- If a pushed commit needs fixing, create a new commit

## Pull Requests

- Use Conventional Commits format for PR title
- Keep description concise

### PR Review Comments

Use Conventional Comments format: `label: subject`

| Label | Purpose |
|-------|---------|
| `praise:` | Highlight something done well |
| `nitpick:` | Minor style/preference, non-blocking |
| `suggestion:` | Propose an alternative approach |
| `issue:` | Something that must be addressed |
| `question:` | Seeking clarification or understanding |
| `thought:` | Share an idea without requiring action |
