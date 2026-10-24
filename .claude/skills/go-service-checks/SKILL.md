---
name: go-service-checks
description: Run formatting, vet and race-detector tests for a Go service under src/ (checkout, product-catalog, ...). Use before committing any Go change.
---

Run these from the service folder (for example `src/checkout`), in order, and stop at the first failure:

1. `gofmt -l .` — must print nothing. If it lists files, run `gofmt -w <file>` on them.
2. `go vet ./...`
3. `go test -race ./...`
4. If you touched a test or suspect flakiness, also run `go test -race -count=20 ./...`.

Report the exact commands you ran and their result. Do not commit while any step fails.

Notes:
- The generated code in `genproto/` is not edited by hand.
- product-catalog loads products from PostgreSQL at runtime; keep new logic in pure functions so it can be unit tested without a database.
