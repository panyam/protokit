# protokit

Shared Go toolkit for building `protoc` plugins. Extracts common patterns — name conversion, field introspection, package resolution, well-known type handling, and test helpers — so each plugin focuses on its own code generation logic.

Used by:
- [`protoc-gen-dal`](https://github.com/panyam/protoc-gen-dal) — DAL layer generation (GORM, Datastore)
- [`mcpkit/ext/protogen`](https://github.com/panyam/mcpkit) — MCP server bindings from proto services

## Install

```bash
go get github.com/panyam/protokit
```

Requires Go 1.24+ and `google.golang.org/protobuf`.

## Packages

### `names` — Name conversion

```go
import "github.com/panyam/protokit/names"

names.ToSnakeCase("GetUserByID")   // "get_user_by_id"
names.ToSnakeCase("HTTPServer")    // "http_server"
names.ToCamelCase("FindBooks")     // "findBooks"
names.ToPascalCase("findBooks")    // "FindBooks"
names.SanitizeIdentifier("user-name") // "user_name"
```

`ToSnakeCase` handles acronyms correctly — `HTTPServer` becomes `http_server`, not `h_t_t_p_server`.

### `fields` — Field introspection

```go
import "github.com/panyam/protokit/fields"

fields.GetFieldKind(field)                    // "string", "int32", "message", "enum"
fields.IsMapField(field)                      // true if proto map<K,V>
fields.GetMapKeyValueFields(field)            // key, value *protogen.Field
fields.IsRepeated(field)                      // true if repeated
fields.IsOptional(field)                      // true if proto3 optional
fields.IsNumericKind("sfixed64")              // true
fields.NormalizeNumericKind("sint32")         // "int32"
```

### `messages` — Message introspection

```go
import "github.com/panyam/protokit/messages"

messages.ExtractPackageName("library.v1.Book")  // "library.v1"
messages.ExtractMessageName("library.v1.Book")  // "Book"
messages.GetFullyQualifiedType(field)            // "library.v1.Book"
messages.IsNestedMessage(msg)                    // true if nested inside another message
messages.GetOneofGroups(msg)                     // ["identity", "content"]
messages.GetBaseFileName("path/to/library.proto") // "library"
messages.BuildMessageIndex(plugin)               // map[string]*protogen.Message across all files
```

### `packages` — Package resolution and imports

```go
import "github.com/panyam/protokit/packages"

// Package info extraction (handles buf ;suffix convention)
info := packages.ExtractPackageInfo(msg) // {ImportPath: "github.com/...", Alias: "v1"}
packages.ExtractGoPackageName(msg)       // "v1"
packages.GetPackageAlias("github.com/example/api/v1") // "v1"

// Path utilities
packages.BuildPackagePath("library.v1")           // "library/v1"
packages.CalculateRelativePath("/a/b", "/a/c")    // "../c"
packages.JoinPaths("./gen", "api", "v1")           // "./gen/api/v1"

// Deduplicating import management
imports := packages.NewImportMap()
imports.Add(packages.ImportSpec{Path: "fmt"})
imports.Add(packages.ImportSpec{Path: "context", Alias: "ctx"})
imports.ToSlice() // sorted, deduplicated
```

### `wellknown` — Well-known type registry

Language-agnostic registry for mapping proto well-known types to target language types. Each plugin registers its own mappings.

```go
import "github.com/panyam/protokit/wellknown"

reg := wellknown.NewRegistry()
reg.Register("google.protobuf.Timestamp", "time.Time", "time", true)
reg.Register("google.protobuf.Duration", "time.Duration", "time", true)

mapping, ok := reg.Get("google.protobuf.Timestamp")
// mapping.TargetType == "time.Time"

reg.IsWellKnown("google.protobuf.Struct") // false (not registered yet)

// Convenience: list of all standard well-known type names
wellknown.WellKnownProtoTypes() // ["google.protobuf.Timestamp", "google.protobuf.Duration", ...]
```

### `testutil` — Test proto descriptor builders

Build in-memory proto descriptors for unit testing protoc plugins without `.proto` files or running `protoc`. Uses `CodeGeneratorRequest` → `protogen.Plugin` — the same path production plugins use.

```go
import "github.com/panyam/protokit/testutil"

plugin := testutil.CreateTestPlugin(t, &testutil.TestProtoSet{
    Files: []testutil.TestFile{{
        Name: "user.proto",
        Pkg:  "user.v1",
        Enums: []testutil.TestEnum{{
            Name:   "Role",
            Values: []testutil.TestEnumValue{{Name: "ADMIN", Number: 0}, {Name: "USER", Number: 1}},
        }},
        Messages: []testutil.TestMessage{
            {
                Name: "GetUserRequest",
                Fields: []testutil.TestField{
                    {Name: "id", Number: 1, TypeName: "string"},
                    {Name: "role", Number: 2, EnumType: "user.v1.Role"},
                },
            },
            {
                Name: "User",
                Fields: []testutil.TestField{
                    {Name: "id", Number: 1, TypeName: "string"},
                    {Name: "name", Number: 2, TypeName: "string"},
                    {Name: "nickname", Number: 3, TypeName: "string", Optional: true},
                    {Name: "tags", Number: 4, TypeName: "string", Repeated: true},
                    {Name: "labels", Number: 5, TypeName: "string", IsMap: true, MapKeyType: "string"},
                },
            },
            {
                Name:   "Notification",
                Oneofs: []string{"channel"},
                Fields: []testutil.TestField{
                    {Name: "email", Number: 1, TypeName: "string", OneofIndex: 0},
                    {Name: "sms", Number: 2, TypeName: "string", OneofIndex: 0},
                },
            },
        },
        Services: []testutil.TestService{{
            Name: "UserService",
            Methods: []testutil.TestMethod{
                {Name: "GetUser", InputType: "user.v1.GetUserRequest", OutputType: "user.v1.User"},
                {Name: "StreamUsers", InputType: "user.v1.GetUserRequest", OutputType: "user.v1.User", ServerStreaming: true},
            },
        }},
    }},
})

// Now use plugin.Files, plugin.Files[0].Messages, etc. in your tests.
```

Supports: all scalar types, enums, maps (any key type), repeated fields, optional fields, oneofs, nested messages, services with unary/streaming methods.

## License

Apache 2.0
