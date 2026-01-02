# features - Go Feature Flags Service

[![CI](https://github.com/serroba/features/actions/workflows/ci.yml/badge.svg)](https://github.com/serroba/features/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/serroba/features/branch/main/graph/badge.svg)](https://codecov.io/gh/serroba/features)
[![Go Report Card](https://goreportcard.com/badge/github.com/serroba/features)](https://goreportcard.com/report/github.com/serroba/features)
[![Go Reference](https://pkg.go.dev/badge/github.com/serroba/features.svg)](https://pkg.go.dev/github.com/serroba/features)

A feature flag evaluation service for Go with rule-based targeting and multi-tenant support.

## Features

- **Rule-Based Targeting** - Evaluate flags based on user attributes, tenant, and custom conditions
- **Multiple Value Types** - Boolean, string, and number flag values
- **Condition Operators** - Equals, not equals, in, not in, exists, starts with
- **Multi-Tenant** - Built-in support for tenant and user context
- **HTTP API** - RESTful API with OpenAPI documentation via Huma

## Quick Start

```bash
go run ./cmd/server --port 8080
```

### Create a Flag

```bash
curl -X POST http://localhost:8080/flags \
  -H "Content-Type: application/json" \
  -d '{
    "key": "dark-mode",
    "type": "bool",
    "enabled": true,
    "defaultValue": {"kind": "bool", "bool": false},
    "rules": [
      {
        "id": "beta-testers",
        "conditions": [{"attr": "plan", "op": "eq", "value": "premium"}],
        "value": {"kind": "bool", "bool": true}
      }
    ]
  }'
```

### Evaluate a Flag

```bash
curl -X POST http://localhost:8080/flags/dark-mode/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "acme-corp",
    "userId": "user-123",
    "attrs": {"plan": "premium"}
  }'
```

## API Endpoints

| Method | Path                      | Description           |
|--------|---------------------------|-----------------------|
| POST   | `/flags`                  | Create a feature flag |
| POST   | `/flags/{key}/evaluate`   | Evaluate a flag       |

## Condition Operators

| Operator      | Description                          |
|---------------|--------------------------------------|
| `eq`          | Attribute equals value               |
| `neq`         | Attribute does not equal value       |
| `in`          | Attribute is in list of values       |
| `not_in`      | Attribute is not in list of values   |
| `exists`      | Attribute exists (is not nil)        |
| `starts_with` | String attribute starts with prefix  |

## Evaluation Order

1. **Disabled Check** - If flag is disabled, return default value with `disabled` reason
2. **Rule Matching** - Evaluate rules in order, first match wins
3. **Default** - If no rules match, return default value

## Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out

# Run linter
golangci-lint run

# Start the server
go run ./cmd/server
```

## License

MIT License - see [LICENSE](LICENSE) for details.
