---
name: testing
description: Write and update tests using the given-when-then pattern for clear test structure. Use when writing unit tests, integration tests, or updating existing test suites.
compatibility:
  - opencode
  - claude
license: MIT
---

# Testing

## Structure

Use the given-when-then pattern:

- **given**: setup and preconditions
- **when**: action being tested
- **then**: expected outcome

Use comments to mark each section when helpful.

## Example

```go
func TestUserService_Create(t *testing.T) {
    // given
    db := setupTestDB(t)
    svc := NewUserService(db)
    
    // when
    user, err := svc.Create("alice@example.com")
    
    // then
    require.NoError(t, err)
    assert.Equal(t, "alice@example.com", user.Email)
}
```
