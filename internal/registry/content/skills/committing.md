---
name: committing
description: Create atomic, well-formatted git commits following Conventional Commits format. Use when committing staged changes, crafting commit messages, or preparing a git commit.
compatibility:
  - opencode
  - claude
license: MIT
---

# Committing

Create a git commit for the staged changes.

## Process

1. **Analyze Changes**
   - Run `git status` to see staged and unstaged files
   - Run `git diff --cached` to review staged changes
   - Understand the purpose and scope of the changes

2. **Craft Commit Message**
   - Use Conventional Commits format: `type(scope): description`
   - Subject line only, max 72 characters
   - Use imperative mood ("add" not "added")
   - Be specific about what changed and why
   - Keep commits focused and atomic (one logical change per commit)
   - Do not add Co-Authored-By lines

3. **Commit Types**
   - `feat`: New feature (triggers release)
   - `fix`: Bug fix (triggers release)
   - `perf`: Performance improvement (triggers release)
   - `docs`: Documentation only
   - `style`: Code style (formatting, semicolons)
   - `refactor`: Code change without feature/fix
   - `test`: Adding or fixing tests
   - `chore`: Maintenance tasks
   - `build`: Build system changes
   - `ci`: CI configuration changes

4. **Scope** (optional)
   - Use lowercase
   - Identify affected component/module
   - Keep it short and consistent

## Breaking Changes

Add `!` after type (or scope) to indicate a breaking change.

For services: API contract changes (endpoints, request/response formats, auth flows)
For libraries: Public API changes (exported functions, types, interfaces)

When uncertain if a change is breaking, ask before committing.

## History Rewriting

- Do not use `git commit --amend` on commits that have been pushed
- If a pushed commit needs fixing, create a new commit

## Examples

Good:
- `feat(auth): add OAuth2 login support`
- `fix(api): handle null response in user endpoint`
- `refactor(utils): simplify date formatting logic`
- `feat(api)!: change response format for /users endpoint`

Bad:
- `updated files` (vague)
- `Fix bug` (no scope, not descriptive)
- `feat: Added new feature for user authentication` (past tense, too long)

## Pull Request Titles

Use the same Conventional Commits format for PR titles.

## Verification

After committing:
- Run `git log -1` to verify the commit
- Ensure the message accurately describes the change
