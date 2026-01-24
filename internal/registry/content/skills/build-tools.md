---
name: build-tools
description: Always use project build tools (Makefile, Mage, npm scripts) instead of running commands directly. Load when running builds, tests, or linting.
compatibility:
  - opencode
  - claude
license: MIT
---

# Build Tools

Before running any build, test, lint, or format command directly, check for project tooling.

## Priority

1. **Check for Makefile** - Run `make help` or read targets
2. **Check for magefile.go** - Run `mage -l` to list targets
3. **Check for package.json scripts** - Run `npm run` to list scripts
4. **Check for project CLI** - Look for `./bin/`, `./scripts/`, or documented tooling

## Why This Matters

- Project tools often have required flags, environment setup, or preprocessing
- Running `go test` directly may miss project-specific test configuration
- Formatters and linters may have project-specific settings

## Common Targets

| Target | Purpose |
|--------|---------|
| `build` | Compile the project |
| `test` | Run tests |
| `lint` | Run linters |
| `fmt` / `format` | Format code |
| `dev` | Start development server |
| `check` | Run all validations |

## Rule

**Never skip project tooling.** If a Makefile exists, use it.
