---
name: pr-review
description: Review pull requests and code for quality, correctness, and maintainability using Conventional Comments. Use when reviewing PRs or code changes.
compatibility:
  - opencode
  - claude
license: MIT
---

# PR Review

Review a pull request systematically and provide actionable feedback.

## Process

1. **Understand Context**
   - Read the PR title and description
   - Understand the purpose and scope of changes
   - Check linked issues or tickets

2. **Review Changes**
   - Use `gh pr diff` or review the diff
   - Consider the overall architecture impact
   - Check for breaking changes

3. **Evaluate Quality**

### Code Quality
- Is the code readable and well-organized?
- Are names meaningful and consistent?
- Is complexity appropriate?

### Correctness
- Does the logic handle edge cases?
- Are error conditions handled?
- Are assumptions documented?

### Testing
- Are tests included for new functionality?
- Do tests cover edge cases?
- Are existing tests still passing?

### Security
- Any hardcoded secrets or credentials?
- Input validation present?
- SQL injection, XSS, or other vulnerabilities?

### Performance
- Any obvious performance issues?
- Unnecessary database queries?
- Memory leaks or resource management issues?

## Feedback Format

Use Conventional Comments format:

```
label (decorations): subject

discussion
```

### Labels

- `praise:` Highlight something done well (leave at least one per review)
- `nitpick:` Trivial preference-based request (non-blocking by nature)
- `suggestion:` Propose an improvement; be explicit on what and why
- `issue:` Specific problem that should be addressed; pair with a suggestion
- `question:` Potential concern needing clarification
- `thought:` Non-blocking idea for consideration or mentoring
- `todo:` Small, trivial, necessary change
- `chore:` Task required before acceptance (link to process docs)
- `note:` Non-blocking; something the reader should be aware of

### Decorations

Add in parentheses after the label:

- `(non-blocking)` - Should not prevent merge
- `(blocking)` - Must be resolved before merge
- `(if-minor)` - Resolve only if the fix is trivial

### Examples

```
praise: Clean separation of concerns here. The service layer is well-defined.

issue (blocking): This SQL query is vulnerable to injection. Use parameterized queries.

suggestion (security): Consider using the framework's built-in sanitizer here.

suggestion (non-blocking): Consider using a Map instead of repeated array lookups.
The O(1) lookup would improve performance for larger datasets.

nitpick: Prefer `const` over `let` for variables that aren't reassigned.

chore: Run the integration tests before merging. See [CI docs](link).
```

## Summary

End with an overall assessment:
- **Approve**: Ready to merge
- **Request Changes**: Issues must be addressed
- **Comment**: Feedback provided, no blocking issues
