# CLAUDE.md — protokit

## Project Context

Shared proto analysis toolkit for building protoc plugins. Extracted from common patterns in `protoc-gen-go-wasmjs` and `protoc-gen-dal`. Provides language-agnostic utilities for proto field introspection, naming conventions, package resolution, well-known type mapping, and test helpers.

## Build & Test

```bash
make build    # go build ./...
make test     # go test -v ./...
make lint     # go vet ./...
make tidy     # go mod tidy
```

## Package Structure

```
protokit/
├── names/       # Naming conventions (camelCase, PascalCase, snake_case)
├── fields/      # Proto field introspection (map, repeated, optional, numeric)
├── messages/    # Message/enum introspection (nested, oneof, fully qualified names)
├── packages/    # Package info extraction, import deduplication, path calculation
├── wellknown/   # Language-agnostic well-known type registry
└── testutil/    # In-memory proto descriptor builders for testing
```

## Key Design Decisions

- **Package-level functions, not struct methods** — all utilities are stateless
- **Language-agnostic** — no hardcoded Go/TS/SQL target types; plugins register their own mappings
- **Sole dependency**: `google.golang.org/protobuf` — no other external deps
- **Generic test helpers** — `testutil.TestMessage.Options` uses `*descriptorpb.MessageOptions` so any plugin can set its own extensions

## Dependencies

- `google.golang.org/protobuf` (protogen, descriptorpb, pluginpb)

## Consumers

- `protoc-gen-go-wasmjs` — WASM bindings + TypeScript client generation
- `protoc-gen-dal` — Data Access Layer code generation (GORM, Datastore)
