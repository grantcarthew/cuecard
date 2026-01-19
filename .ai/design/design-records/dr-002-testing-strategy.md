# dr-002: Testing Strategy

- Date: 2026-01-19
- Status: Accepted
- Category: testing

## Problem

Without a defined testing strategy, code quality depends on ad-hoc testing. GUI applications have components with varying testability - some logic is easily unit tested while UI behavior is harder to automate. Cuecard needs a clear approach to testing that balances coverage with practical constraints.

## Decision

Use Go's standard `testing` package with table-driven tests. Focus testing effort on high-value, easily testable packages. Accept manual testing for UI components where Fyne's test support is limited.

Testing by package:

| Package | Strategy | Coverage Goal |
|---------|----------|---------------|
| `internal/config` | Unit tests for CUE parsing, validation | High |
| `internal/prompt` | Unit tests for frontmatter parsing, variable substitution, filename generation | High |
| `internal/clipboard` | Integration tests with interface abstraction | Medium |
| `internal/watcher` | Integration tests with temp directories | Medium |
| `internal/ui` | Manual testing; Fyne test utilities where practical | Low |

Test execution:

- `go test ./...` runs all tests
- Tests run in CI on every push

## Why

Go's built-in testing:

- Sufficient for project needs without external dependencies
- Well-documented and familiar to Go developers
- Integrated with go toolchain (go test, go cover)

Table-driven tests:

- Idiomatic Go pattern
- Easy to add edge cases
- Tests serve as documentation of expected behavior

Focus on pure logic packages:

- Config and prompt parsing is where bugs are most likely
- These packages have no external dependencies, easy to test
- Variable substitution has many edge cases worth covering

Limited UI testing investment:

- Fyne's testing support is basic
- UI tests tend to be fragile and slow
- Manual testing before releases is acceptable for v1

## Trade-offs

Accept:

- UI bugs may slip through without automated UI tests
- Manual testing required before releases
- No mocking framework (use interfaces instead)

Gain:

- Simple test setup, no external dependencies
- Fast test execution
- High confidence in core parsing/substitution logic
- Tests serve as documentation

## Alternatives

testify/assert:

- Pro: Nicer assertion syntax
- Pro: Mock generation with testify/mock
- Con: External dependency to maintain
- Con: Adds complexity for marginal benefit
- Rejected: Standard library sufficient for this project size

Fyne test automation:

- Pro: Automated UI testing possible
- Pro: Catch visual regressions
- Con: Complex setup and maintenance
- Con: Fragile tests that break with UI changes
- Con: Limited Fyne test utilities available
- Rejected: ROI too low for v1

End-to-end testing:

- Pro: Tests full user workflows
- Con: Slow execution
- Con: Complex setup with GUI automation
- Con: Flaky tests common
- Rejected: Not practical for desktop GUI app

## Structure

Test file organization:

- Test files adjacent to source: `config.go` → `config_test.go`
- Test data in `testdata/` subdirectories where needed
- Table-driven test format for multiple cases

Example test structure:

```
internal/
  config/
    config.go
    config_test.go
    testdata/
      valid.cue
      invalid.cue
  prompt/
    parser.go
    parser_test.go
    substitution.go
    substitution_test.go
    testdata/
      simple.md
      with-variables.md
```

## Validation

Tests pass when:

- `go test ./...` exits with code 0
- No race conditions detected with `go test -race ./...`
- Coverage meets minimum threshold per package (config, prompt: 80%+)
