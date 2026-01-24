---
name: releasing
description: Drafts release notes from commits, proposes semantic version bumps, and generates changelog entries. Use when preparing a tagged release, creating a new version, or generating release notes.
compatibility:
  - opencode
  - claude
license: MIT
---

# Releasing

Prepares tagged releases with proper versioning and documentation.

## Capabilities

- Drafts release notes from merged PRs since the last tag
- Proposes a semantic version bump based on commit types
- Generates a changelog entry
- Provides a copy-pasteable `gh release create` command

## Process

1. Analyze commits since the last tag
2. Suggest a version number based on conventional commits
3. Draft release notes for review and editing

## Version Bump Rules

- `feat:` commits trigger a minor version bump
- `fix:` commits trigger a patch version bump
- `BREAKING CHANGE:` or `!` trigger a major version bump
