---
name: build-tools
description: Build tool preferences for Makefile, Mage, and project CLI tooling. Load when running builds, tests, linting, or using project tooling.
compatibility:
  - opencode
  - claude
license: MIT
---

## Build Tools & CLI

- Prefer CLI tools over manual code manipulation
- Use Makefile if project has one
- Use Mage if project has magefile.go
- Common targets: build, test, lint, fmt
- Use project linters and formatters when available
- Run tests through project tooling, not manually
