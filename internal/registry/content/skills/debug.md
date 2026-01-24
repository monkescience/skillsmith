---
name: debug
description: Debug issues systematically using a structured approach. Use when troubleshooting bugs, investigating errors, diagnosing problems, or fixing unexpected behavior.
compatibility:
  - opencode
  - claude
license: MIT
---

# Debugging

A systematic approach to finding and fixing bugs.

## 1. Reproduce the Issue

- Get exact steps to reproduce the bug
- Note the expected vs actual behavior
- Identify the minimal reproduction case

## 2. Gather Information

- Read error messages carefully
- Check logs for relevant entries
- Identify when the issue started (recent changes?)
- Note any patterns (specific inputs, timing, environment)

## 3. Form Hypotheses

Based on the symptoms, list possible causes:
- Recent code changes
- Configuration issues
- Data corruption
- Race conditions
- External dependencies

## 4. Test Hypotheses

For each hypothesis:
1. Predict what you should observe if it's correct
2. Design a test to verify
3. Execute the test
4. Evaluate results

## 5. Isolate the Problem

Use binary search to narrow down:
- Comment out code sections
- Add logging at key points
- Use a debugger to step through
- Simplify inputs

## 6. Fix and Verify

- Make the minimal fix
- Add a test that would have caught this bug
- Verify the fix in the original context
- Check for similar issues elsewhere

## 7. Document

- Add comments if the fix isn't obvious
- Update documentation if needed
- Share learnings with the team
