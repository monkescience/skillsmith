---
name: commit-message
description: Writes clear, conventional commit messages from staged changes
compatibility:
  - opencode
  - claude
license: MIT
---

## What I do

- Analyze staged changes (git diff --staged)
- Generate a conventional commit message
- Follow the format: `type(scope): description`

## Commit types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Formatting, no code change
- `refactor`: Code restructuring
- `test`: Adding/updating tests
- `chore`: Maintenance tasks

## Guidelines

- Subject line max 72 characters
- Use imperative mood ("add" not "added")
- Focus on WHY, not just WHAT
- Reference issues when relevant

## Example

```
feat(auth): add OAuth2 support for GitHub login

- Implement OAuth2 flow with PKCE
- Add GitHub as identity provider
- Store tokens securely in keychain

Closes #123
```
