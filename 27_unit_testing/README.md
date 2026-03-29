# Unit Test Convention

This convention provides a **consistent, scalable, and maintainable** approach to unit testing in Go.

## 0. File Structure

The example files for this module are organized as follows:
- [domain/user.go](domain/user.go): Domain entity and interface definitions.
- [service/user_service.go](service/user_service.go): Service implementation.
- [service/user_service_test.go](service/user_service_test.go): **The Unit Test example code.**
- [mocks/mocks.go](mocks/mocks.go): Mock implementations using `testify/mock`.

---

## 1. Naming Convention

### Function Name
Use the format `TestServiceName_MethodName(t *testing.T)`.

**Format:**
```go
func TestServiceName_MethodName(t *testing.T)
```

**Examples:**
- `func TestUserService_GetByID(t *testing.T)`
- `func TestOrderService_CreateOrder(t *testing.T)`
- `func TestPaymentService_ProcessRefund(t *testing.T)`

**Why this format?**
- **Simple and consistent**: directly shows the service and method being tested.
- **Easy to filter**: when running tests with `go test -run TestUserService_`.
- **Scenario details**: are in the description field within test cases, not in the function name.
- **Common practice**: in the Go community for service layer testing.

### Test Case Description
Use the format `"STATUS: Description"`.
- `SUCCESS`: happy path scenarios.
- `ERROR`: error scenarios.
- `FAILED`: business logic failures.

---

## 2. Structure

A standard unit test in this convention uses:
- **Table-driven tests**: to handle multiple scenarios efficiently.
- **Mocker struct**: to centralize and pass mock dependencies.
- **mockSetup function**: to allow isolated mock configuration per test case.

### Complete Example
```go
func TestUserService_GetByID(t *testing.T) {
    type Mocker struct {
        repo         *mocks.UserRepository
        mockProducer *mocks.EventPublisher
    }

    testCases := []struct {
        desc      string
        mockSetup func(m *Mocker)
        wantErr   bool
        expected  *domain.User
    }{
        {
            desc:    "SUCCESS: Get User by ID",
            wantErr: false,
            mockSetup: func(m *Mocker) {
                m.repo.On(
                    "GetByID",
                    mock.Anything,
                    mock.AnythingOfType("int64"),
                ).Return(&domain.User{
                    ID:    100,
                    Name:  "John Doe",
                    Email: "dummy@mail.co",
                }, nil)
            },
            expected: &domain.User{
                ID:    100,
                Name:  "John Doe",
                Email: "dummy@mail.co",
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.desc, func(t *testing.T) {
            // 1. Initialize fresh mocks per test case
            repo := new(mocks.UserRepository)
            producer := new(mocks.EventPublisher)
            
            mocker := &Mocker{
                repo:         repo,
                mockProducer: producer,
            }

            // 2. Setup mock expectations
            tc.mockSetup(mocker)

            // 3. Initialize service with mocks
            service := NewService(repo, producer)

            // 4. Execute method
            actual, err := service.GetByID(context.Background(), 100)

            // 5. Assert results
            if tc.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.EqualValues(t, tc.expected, actual)
            }

            // 6. Verify mock expectations
            repo.AssertExpectations(t)
            producer.AssertExpectations(t)
        })
    }
}
```

---

## 3. Mocking with Mockery

### Setup Mockery
Install mockery:
```bash
go install github.com/vektra/mockery/v2@latest
```

Generate mocks from interfaces:
```bash
# Generate specific interface
mockery --name=UserRepository --output=./internal/domain --output=./mocks --outpkg=mocks

# Or generate all interfaces at once
mockery --all
```

**Why `with-expecter: false`?**
- Uses the more explicit and familiar `.On()` syntax.
- Aligns with the pattern in the examples.
- More flexible for complex parameter matching.
- Consistent with `testify/mock` documentation.

### Mock Parameter Matchers
1.  **`mock.Anything`**: Match any value (useful for Context, IDs, etc.).
2.  **`mock.AnythingOfType("type")`**: Type-safe matching (recommended for most cases).
3.  **Specific Value**: When you need exact matching for edge cases or specific IDs.
4.  **`mock.MatchedBy`**: Custom matcher function for complex validation logic.

| Matcher | Use Case | Example |
| :--- | :--- | :--- |
| `mock.Anything` | Parameters not relevant to the test | Context, metadata, trace ID |
| `mock.AnythingOfType("T")` | Ensure type safety without checking specific values | **Recommended for most cases** |
| Specific Value | Test cases needing exact matching | Edge cases, specific IDs |

---

## 4. Mock Return Values

### Examples
```go
// Success case - return data + nil error
m.repo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
    Return(&domain.User{ID: 100, Name: "John"}, nil)

// Error case - return nil + error
m.repo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
    Return(nil, assert.AnError)

// Specific error
m.repo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
    Return(nil, errors.New("user not found"))
```

### Why `assert.AnError`?
- **Generic error**: for testing error handling behavior.
- **No specific messages**: no need to create unique error strings for every test.
- **Behavior focus**: focus on "error vs no error", not content.
- **Easy maintenance**: simplifies test updates.

---

## 5. Running Tests

### Basic Commands
- Run all tests: `go test ./...`
- Run with coverage: `go test ./... -cover`
- Run specific package: `go test ./internal/service/...`
- Run specific test function: `go test -run TestUserService_GetByID ./...`
- Run with verbose output: `go test -v ./...`

### Generate Coverage Reports
```bash
# Generate HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Generate terminal coverage report
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

---

## 6. Coverage Threshold

### Minimum Acceptable: 80%
80% coverage is the "sweet spot" between quality and effort. Beyond 80% often has diminishing returns.

### Check Threshold via Command
```bash
go test ./... -coverprofile=coverage.out && \
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | \
awk '{if ($1 < 80) exit 1}'
```

### Coverage Exceptions
The following should be excluded from coverage requirements:
- **Generated code** (protobuf, mocks)
- **Main functions and CLI code** (hard to test, usually excluded)
- **Trivial getters/setters** (no logic, no need to test)

---

## 7. PR Requirements
- **All tests must pass**: no failing tests allowed.
- **Coverage 80%**: minimum threshold must be met.
- **No skipped tests**: no test files in the diff should be skipped.
- **Mock expectations asserted**: all mocks must have `AssertExpectations(t)` called.

## Summary
Following these conventions will ensure high-quality, reliable, and maintainable tests across your codebase:
- Clear naming conventions.
- Table-driven patterns for comprehensive coverage.
- Proper mocking with mockery.
- Consistent assertion patterns.
- Test isolation and independence.
