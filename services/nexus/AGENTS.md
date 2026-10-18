# Nexus service guidance

Nexus is the backend boundary for application services that read Grist data. Keep Grist credentials and document IDs on the server.

## Request flow

- Keep Expenses HTTP parsing and response writing in `internal/handler/`.
- Put Expenses request validation and application behavior in `internal/service/`.
- Keep Grist table names, column names, and SQL in `internal/repository/`.
- Keep HTTP communication with Grist in `internal/gristapi/`.
- Keep shared request and response types in `internal/model/`.
- Wire handlers to services and repositories in `main.go`.
- New application endpoints should follow `handler -> service -> repository -> grist.Client`.

## API behavior

- Expense endpoints are read-only and live under `/expenses/v1`.
- `GET /expenses/v1/catalog` returns account, category, and tag choices.
- `GET /expenses/v1/monthly-summary` accepts `month` and CSV `accounts`, `categories`, and `tags` query parameters.
- These endpoints rely on private tailnet reachability and do not use application-level authentication.
- Return generic upstream errors to callers. Do not expose Grist errors, credentials, or private records in logs.
- Preserve the JSON response shapes consumed by the Expenses service.

## Go style

- Keep code idiomatic and formatted with `gofmt`.
- Use request contexts for upstream calls.
- Avoid new dependencies unless the implementation requires them.
