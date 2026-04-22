# Repository Guidelines

## Project Structure & Module Organization
`main.go` boots the Cobra CLI. `cmd/` defines commands such as `transfer`, `config/` loads Viper settings and MySQL connections, `services/aryan` and `services/fararavand` contain ERP-specific clients, and `sync_data/` holds sync flows by domain (`invoice.go`, `product.go`, `customer.go`). Use `repository/` and `repository/database/` for persistence and sync checkpoints, `models/` for shared data structures, and `utility/logger/` for reusable logging.

## Build, Test, and Development Commands
`go run . transfer -c env.yml` runs the sync job with an explicit config file.
`go build -o bin/erp-job .` builds the CLI binary for local execution.
`go test ./...` checks that all packages compile; the repository currently has no `_test.go` files, so new work should add coverage.
`go vet ./...` catches common Go issues before review.
`gofmt -w ./...` or `gofmt -w path/to/file.go` should be run before committing.

## Coding Style & Naming Conventions
Follow standard Go formatting and imports; let `gofmt` decide spacing and tabs. Keep package names lowercase, exported identifiers in `PascalCase`, and internal helpers in `camelCase`. Name files by business area or responsibility, for example `invoice.go`, `customer.go`, or `database.go`. Prefer small functions that handle one sync step and inject dependencies through constructors.

## Testing Guidelines
Place tests beside the code they cover using the `_test.go` suffix. Prefer table-driven tests for model mapping, repository behavior, and sync edge cases. When touching ERP integrations, move transform logic into testable helpers so behavior can be validated without live APIs or MySQL. Run `go test ./...` before opening a PR.

## Commit & Pull Request Guidelines
Recent history favors short imperative commit subjects, often with emoji prefixes like `:zap:` and `:pencil:`. Keep commits focused on one logical change. Pull requests should state which sync flow is affected, list config or schema changes, link the related issue, and include sample payloads or logs when behavior changes.

## Security & Configuration Tips
Do not commit real ERP credentials or database passwords. The app expects a local `env.yml` by default; if you add new config keys, document them in the PR description and keep example values sanitized.
