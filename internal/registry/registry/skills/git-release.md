---
name: git-release
description: Creates consistent releases with changelogs and version bumps
category: workflow
compatibility:
  - opencode
  - claude
tags:
  - git
  - release
  - changelog
author: skillsmith
license: MIT
---

## What I do

- Draft release notes from merged PRs since the last tag
- Propose a semantic version bump based on commit types
- Generate a changelog entry
- Provide a copy-pasteable `gh release create` command

## When to use me

Use this skill when you are preparing a tagged release.

## How to use

1. Ask me to prepare a release
2. I'll analyze commits since the last tag
3. I'll suggest a version number based on conventional commits
4. I'll draft release notes you can review and edit

## Version bump rules

- `feat:` commits trigger a minor version bump
- `fix:` commits trigger a patch version bump
- `BREAKING CHANGE:` or `!` trigger a major version bump
