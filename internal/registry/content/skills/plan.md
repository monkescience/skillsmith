---
name: plan
description: Create detailed implementation plans before writing code. Use when planning new features, complex tasks, architectural changes, or multi-step implementations.
compatibility:
  - opencode
  - claude
license: MIT
---

# Plan

Create a detailed implementation plan before writing code.

## Planning Process

1. **Understand Requirements**
   - What is the goal?
   - Who are the users/consumers?
   - What are the acceptance criteria?

2. **Research Existing Code**
   - Find similar patterns in the codebase
   - Identify code to reuse or extend
   - Note conventions and standards used

3. **Identify Scope**
   - List all files that need changes
   - Identify new files to create
   - Consider dependencies and imports

## Plan Structure

### Overview
Brief description of what will be implemented and why.

### Approach
High-level strategy and key decisions:
- Architecture choices
- Design patterns to use
- Trade-offs considered

### Implementation Steps
Ordered list of concrete steps:
1. Step description
   - Files affected
   - Key changes
2. Next step...

### Testing Strategy
- Unit tests needed
- Integration tests
- Edge cases to cover

### Risks and Considerations
- Potential issues
- Dependencies on other work
- Migration needs

## Guidelines

- **Be specific**: Name files, functions, and types
- **Stay focused**: Only what's needed for this task
- **Consider existing patterns**: Follow codebase conventions
- **Think about errors**: How will failures be handled?
- **Keep it simple**: Avoid over-engineering

## Example Plan

```markdown
## Overview
Add user authentication via OAuth2.

## Approach
Use existing session middleware. Implement OAuth2 flow
with provider abstraction for future providers.

## Steps
1. Add OAuth2 config to settings
   - config/settings.go: Add OAuth2Config struct
2. Create auth provider interface
   - internal/auth/provider.go: Provider interface
3. Implement Google provider
   - internal/auth/google.go: GoogleProvider
4. Add auth routes
   - internal/routes/auth.go: /login, /callback
5. Update middleware
   - internal/middleware/session.go: Check auth state

## Testing
- Unit: Provider interface mocks
- Integration: OAuth flow with test server
- E2E: Full login flow

## Risks
- Token refresh handling
- Session expiry edge cases
```
