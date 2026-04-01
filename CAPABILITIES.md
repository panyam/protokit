# protokit

## Version
0.0.1

## Provides
- name-conversion: ToCamelCase, ToPascalCase, ToSnakeCase, SanitizeIdentifier for proto naming
- field-introspection: IsMapField, GetMapKeyValueFields, IsRepeated, IsOptional, IsNumericKind, GetFieldKind
- message-introspection: IsNestedMessage, ExtractPackageName, ExtractMessageName, GetOneofGroups, BuildMessageIndex
- package-resolution: ExtractPackageInfo, GetPackageAlias, ImportSpec/ImportMap with deduplication
- path-calculation: CalculateRelativePath, BuildPackagePath, NormalizePath, JoinPaths
- wellknown-registry: Language-agnostic well-known type registry with plugin-specific target mappings
- test-proto-builders: In-memory proto descriptor construction for unit testing protoc plugins

## Module
github.com/panyam/protokit

## Location
newstack/protokit/main

## Stack Dependencies
- None

## Integration

### Go Module
```go
// go.mod
require github.com/panyam/protokit v0.0.1

// Local development
replace github.com/panyam/protokit => ~/newstack/protokit/main
```

### Key Imports
```go
import "github.com/panyam/protokit/names"
import "github.com/panyam/protokit/fields"
import "github.com/panyam/protokit/messages"
import "github.com/panyam/protokit/packages"
import "github.com/panyam/protokit/wellknown"
import "github.com/panyam/protokit/testutil"
```

## Status
Active

## Conventions
- Package-level functions (no unnecessary structs for stateless operations)
- Language-agnostic core (no Go/TS/SQL target types hardcoded)
- Plugins register their own well-known type mappings
- testutil uses generic Options fields (MessageOptions, FieldOptions) for plugin-specific annotations

## Migrations
