# Practice 4: Unit Testing (Practice 3)

Unit testing untuk project [Practice 3: Net/HTTP with Separation of Concerns](../practice_03_nethttp_soc/README.md), mengikuti konvensi dari [Materi 27: Unit Testing](../27_unit_testing/README.md).

## Cara Jalankan

```bash
# Jalankan semua test
cd practice_04_unit_testing
go test ./...

# Jalankan dengan verbose output
go test -v ./...

# Jalankan dengan coverage
go test -v -cover ./...

# Jalankan test spesifik
go test -v -run TestPersonService_CreatePerson ./service/...
go test -v -run TestPersonHandler_GetPerson ./handler/...
```

## Struktur File

```
practice_04_unit_testing/
├── go.mod
├── entity/
│   └── person.go                    # Domain entity
├── dto/
│   ├── person_dto.go                # Request/Response DTO
│   └── base_response.go            # Base response wrapper
├── repository/
│   └── person_repository.go        # Repository interface + implementation
├── service/
│   ├── person_service.go           # Service implementation
│   └── person_service_test.go      # ✅ Unit test service layer
├── handler/
│   ├── person_handler.go           # HTTP handler
│   ├── person_route.go             # Route mapping
│   └── person_handler_test.go      # ✅ Unit test handler layer
├── server/
│   ├── helper.go                   # API path helper
│   ├── http.go                     # Middleware
│   └── server_test.go              # ✅ Unit test server utilities
└── mocks/
    ├── person_repository_mock.go   # Mock PersonRepository
    └── person_service_mock.go      # Mock PersonService
```

## Coverage Results

| Package | Coverage |
|---------|----------|
| `service` | **100.0%** |
| `handler` | **92.1%** |
| `server` | **100.0%** |

## Yang Ditest

### Service Layer (`service/person_service_test.go`)
- `TestPersonService_CreatePerson`: SUCCESS + ERROR
- `TestPersonService_GetPerson`: SUCCESS + ERROR
- `TestPersonService_GetAllPersons`: SUCCESS + Empty List + ERROR

### Handler Layer (`handler/person_handler_test.go`)
- `TestPersonHandler_CreatePerson`: SUCCESS + Service Error + Invalid JSON Body
- `TestPersonHandler_GetAllPersons`: SUCCESS + Service Error
- `TestPersonHandler_GetPerson`: SUCCESS + Not Found + Invalid ID

### Server Utilities (`server/server_test.go`)
- `TestNewAPIPath`: GET, POST, GET with ID
- `TestApplicationMiddlewareResponse`: Content-Type header
- `TestHandleRoutesNotFound`: Route exists + Route not found

## Konvensi yang Digunakan
- ✅ Naming: `TestServiceName_MethodName`
- ✅ Table-driven tests dengan `Mocker` struct
- ✅ `mockSetup` function per test case
- ✅ Fresh mocks per test case (isolation)
- ✅ `assert.AnError` untuk generic error testing
- ✅ `AssertExpectations(t)` di setiap test case
- ✅ Coverage minimum 80%
