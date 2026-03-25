# Belajar Go

Materi dasar Golang.

## Cara Jalankan

```bash
go run main.go       # jalankan semua materi
go run main.go 1     # jalankan materi tertentu (1-24)
```

## Daftar Materi

| No | Topik | File |
|----|-------|------|
| 1 | [Variables & Data Types](01_variables/README.md) | `01_variables/variables.go` |
| 2 | [Constants, Iota & Operators](02_constants/README.md) | `02_constants/constants.go` |
| 3 | [Functions](03_functions/README.md) | `03_functions/functions.go` |
| 4 | [If/Else & Switch](04_conditions/README.md) | `04_conditions/conditions.go` |
| 5 | [Looping](05_looping/README.md) | `05_looping/looping.go` |
| 6 | [Basic Error Handling](06_errors/README.md) | `06_errors/errors.go` |
| 7 | [Structs & Object Modeling](07_structs/README.md) | `07_structs/structs.go` |
| 8 | [Methods (Value vs Pointer receivers)](08_methods/README.md) | `08_methods/methods.go` |
| 9 | [Pointer Concepts](09_pointers/README.md) | `09_pointers/pointers.go` |
| 10 | [Interfaces as Contracts](10_interfaces/README.md) | `10_interfaces/interfaces.go` |
| 11 | [Dependency Injection](11_dependency_injection/README.md) | `11_dependency_injection/di.go` |
| 12 | [Package System](12_packages/README.md) | `12_packages/packages.go` |
| 13 | [Exported vs Unexported](13_exported/README.md) | `13_exported/exported.go` |
| 14 | [go mod tidy](14_go_mod_tidy/README.md) | `14_go_mod_tidy/gomodtidy.go` |
| 15 | [Separation of Concerns](15_separation_of_concerns/README.md) | `15_separation_of_concerns/separation.go` |
| 16 | [Error Wrapping](16_error_wrapping/README.md) | `16_error_wrapping/errorwrapping.go` |
| 17 | [Database Integration](17_database/README.md) | `17_database/database.go` |
| 18 | [Basic HTTP Server](18_http_server/README.md) | `18_http_server/httpserver.go` |
| 19 | [Handlers](19_handlers/README.md) | `19_handlers/handlers.go` |
| 20 | [JSON Encoding/Decoding](20_json/README.md) | `20_json/json.go` |
| 21 | [HTTP Methods](21_http_methods/README.md) | `21_http_methods/httpmethods.go` |
| 22 | [Status Codes](22_status_codes/README.md) | `22_status_codes/statuscodes.go` |
| 23 | [Middleware Concepts](23_middleware/README.md) | `23_middleware/middleware.go` |

## Practice Projects

| No | Topik | Folder |
|----|-------|--------|
| P1 | [Practice 1: Simple API Fetch (Integrated)](practice_01_api_fetch/README.md) | `practice_01_api_fetch/` |
| P2 | [Practice 2: PostgreSQL CRUD API (Standalone)](practice_02_postgres_crud/README.md) | `practice_02_postgres_crud/` |
| P3 | [Practice 3: Net/HTTP with Separation of Concerns](practice_03_nethttp_soc/README.md) | `practice_03_nethttp_soc/` |

## Challenges

| No | Topik | Folder |
|----|-------|--------|
| C1 | [Challenge 1: Currency Conversion](challenge_1/README.md) | `challenge_1/` |
| C2 | [Challenge 2: Card Identification](challenge_2/README.md) | `challenge_2/` |
| C3 | [Challenge 3: REST API Bank (SOC)](challenge_3/README.md) | `challenge_3/` |

## Project Structure

```
belajar-go/
├── go.mod
├── main.go
├── 01_variables/
├── 02_constants/
├── 03_functions/
├── 04_conditions/
├── 05_looping/
├── 06_errors/
├── 07_structs/
├── 08_methods/
├── 09_pointers/
├── 10_interfaces/
├── 11_dependency_injection/
├── 12_packages/
│   ├── mathutil/
│   └── stringutil/
├── 13_exported/
├── 14_go_mod_tidy/
├── 15_separation_of_concerns/
├── 16_error_wrapping/
├── 17_database/
├── 18_http_server/
├── 19_handlers/
├── 20_json/
├── 21_http_methods/
├── 22_status_codes/
├── 23_middleware/
├── practice_01_api_fetch/
├── practice_02_postgres_crud/  ← standalone project
├── practice_03_nethttp_soc/    ← net/http SoC project
├── challenge_1/
├── challenge_2/
└── challenge_3/                ← REST API Bank (SOC)
```

Setiap folder berisi file `.go` (kode) dan `README.md` (penjelasan).
