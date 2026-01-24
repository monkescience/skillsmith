---
name: refactoring
description: Refactor code safely while preserving behavior. Use when improving code structure, reducing duplication, renaming for clarity, extracting functions, or reorganizing modules.
compatibility:
  - opencode
  - claude
license: MIT
---

# Refactoring

When refactoring code, follow these steps:

## 1. Understand First

- Read and understand the existing code thoroughly
- Identify the current behavior and any edge cases
- Note existing tests that verify the behavior

## 2. Plan the Refactoring

- Identify the specific improvements to make
- Break down into small, incremental changes
- Each change should keep tests passing

## 3. Common Refactorings

- **Extract Function**: Move code into a well-named function
- **Inline Function**: Replace function call with its body when it adds no clarity
- **Rename**: Improve names of variables, functions, or types
- **Extract Variable**: Name intermediate results for clarity
- **Move**: Relocate code to a more appropriate module
- **Remove Duplication**: Consolidate repeated code patterns

## 4. Safety Checklist

- [ ] Tests pass before starting
- [ ] Make one change at a time
- [ ] Run tests after each change
- [ ] Commit working states frequently
- [ ] No behavior changes - only structure

## 5. When to Stop

Stop refactoring when:
- The code clearly expresses its intent
- Tests are comprehensive and passing
- The change achieves its goal without scope creep
