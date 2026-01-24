---
name: build-tools
description: Preferences for build tools and CLI usage - Makefile, Mage, project tooling
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
