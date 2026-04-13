package testutil

import (
	"testing"

	"github.com/panyam/protokit/fields"
	"github.com/panyam/protokit/messages"
)

func TestCreateTestPlugin(t *testing.T) {
	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name: "library.proto",
				Pkg:  "library.v1",
				Messages: []TestMessage{
					{
						Name: "Book",
						Fields: []TestField{
							{Name: "title", Number: 1, TypeName: "string"},
							{Name: "page_count", Number: 2, TypeName: "int32"},
							{Name: "available", Number: 3, TypeName: "bool"},
						},
					},
					{
						Name: "Author",
						Fields: []TestField{
							{Name: "name", Number: 1, TypeName: "string"},
							{Name: "books", Number: 2, TypeName: "library.v1.Book", Repeated: true},
						},
					},
				},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	if plugin == nil {
		t.Fatal("Expected non-nil plugin")
	}

	if len(plugin.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(plugin.Files))
	}

	file := plugin.Files[0]
	if len(file.Messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(file.Messages))
	}

	book := file.Messages[0]
	if string(book.Desc.Name()) != "Book" {
		t.Errorf("Expected message name 'Book', got '%s'", book.Desc.Name())
	}
	if len(book.Fields) != 3 {
		t.Errorf("Expected 3 fields on Book, got %d", len(book.Fields))
	}
}

func TestMapField(t *testing.T) {
	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name: "test.proto",
				Pkg:  "test.v1",
				Messages: []TestMessage{
					{
						Name: "Config",
						Fields: []TestField{
							{Name: "settings", Number: 1, TypeName: "string", IsMap: true, MapKeyType: "string"},
						},
					},
				},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	file := plugin.Files[0]
	config := file.Messages[0]

	if len(config.Fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(config.Fields))
	}

	settingsField := config.Fields[0]
	if !fields.IsMapField(settingsField) {
		t.Error("Expected settings to be a map field")
	}

	keyField, valueField := fields.GetMapKeyValueFields(settingsField)
	if keyField == nil || valueField == nil {
		t.Fatal("Expected non-nil key and value fields")
	}
	if fields.GetFieldKind(keyField) != "string" {
		t.Errorf("Expected key kind 'string', got '%s'", fields.GetFieldKind(keyField))
	}
	if fields.GetFieldKind(valueField) != "string" {
		t.Errorf("Expected value kind 'string', got '%s'", fields.GetFieldKind(valueField))
	}
}

func TestBuildMessageIndex(t *testing.T) {
	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name: "a.proto",
				Pkg:  "pkg.a",
				Messages: []TestMessage{
					{Name: "Foo", Fields: []TestField{{Name: "x", Number: 1, TypeName: "string"}}},
				},
			},
			{
				Name: "b.proto",
				Pkg:  "pkg.b",
				Messages: []TestMessage{
					{Name: "Bar", Fields: []TestField{{Name: "y", Number: 1, TypeName: "int32"}}},
				},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	index := messages.BuildMessageIndex(plugin)

	if _, ok := index["pkg.a.Foo"]; !ok {
		t.Error("Expected 'pkg.a.Foo' in message index")
	}
	if _, ok := index["pkg.b.Bar"]; !ok {
		t.Error("Expected 'pkg.b.Bar' in message index")
	}
}

func TestRepeatedField(t *testing.T) {
	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name: "test.proto",
				Pkg:  "test.v1",
				Messages: []TestMessage{
					{
						Name: "List",
						Fields: []TestField{
							{Name: "items", Number: 1, TypeName: "string", Repeated: true},
							{Name: "count", Number: 2, TypeName: "int32"},
						},
					},
				},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	list := plugin.Files[0].Messages[0]

	if !fields.IsRepeated(list.Fields[0]) {
		t.Error("Expected 'items' to be repeated")
	}
	if fields.IsRepeated(list.Fields[1]) {
		t.Error("Expected 'count' to not be repeated")
	}
}

func TestEnumField(t *testing.T) {
	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name: "test.proto",
				Pkg:  "test.v1",
				Enums: []TestEnum{
					{
						Name: "Status",
						Values: []TestEnumValue{
							{Name: "UNKNOWN", Number: 0},
							{Name: "ACTIVE", Number: 1},
							{Name: "INACTIVE", Number: 2},
						},
					},
				},
				Messages: []TestMessage{
					{
						Name: "Msg",
						Fields: []TestField{
							{Name: "status", Number: 1, EnumType: "test.v1.Status"},
						},
					},
				},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	msg := plugin.Files[0].Messages[0]

	if len(msg.Fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(msg.Fields))
	}
	if msg.Fields[0].Desc.Kind().String() != "enum" {
		t.Errorf("Expected enum kind, got %s", msg.Fields[0].Desc.Kind())
	}
	enumDesc := msg.Fields[0].Desc.Enum()
	if enumDesc == nil {
		t.Fatal("Expected non-nil enum descriptor")
	}
	if enumDesc.Values().Len() != 3 {
		t.Errorf("Expected 3 enum values, got %d", enumDesc.Values().Len())
	}
}

func TestOneofField(t *testing.T) {
	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name: "test.proto",
				Pkg:  "test.v1",
				Messages: []TestMessage{
					{
						Name:   "Msg",
						Oneofs: []string{"value"},
						Fields: []TestField{
							{Name: "text", Number: 1, TypeName: "string", OneofIndex: 0},
							{Name: "number", Number: 2, TypeName: "int32", OneofIndex: 0},
						},
					},
				},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	msg := plugin.Files[0].Messages[0]

	if len(msg.Oneofs) != 1 {
		t.Fatalf("Expected 1 oneof, got %d", len(msg.Oneofs))
	}
	if string(msg.Oneofs[0].Desc.Name()) != "value" {
		t.Errorf("Expected oneof name 'value', got '%s'", msg.Oneofs[0].Desc.Name())
	}

	// Both fields should belong to the oneof.
	for _, f := range msg.Fields {
		if f.Desc.ContainingOneof() == nil {
			t.Errorf("Field %s should be in a oneof", f.Desc.Name())
		}
	}
}

func TestOptionalField(t *testing.T) {
	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name: "test.proto",
				Pkg:  "test.v1",
				Messages: []TestMessage{
					{
						Name: "Msg",
						Fields: []TestField{
							{Name: "required", Number: 1, TypeName: "string"},
							{Name: "optional", Number: 2, TypeName: "string", Optional: true},
						},
					},
				},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	msg := plugin.Files[0].Messages[0]

	if msg.Fields[0].Desc.HasOptionalKeyword() {
		t.Error("'required' field should not have optional keyword")
	}
	if !msg.Fields[1].Desc.HasOptionalKeyword() {
		t.Error("'optional' field should have optional keyword")
	}
}

func TestServiceSupport(t *testing.T) {
	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name: "test.proto",
				Pkg:  "test.v1",
				Messages: []TestMessage{
					{Name: "GetReq", Fields: []TestField{{Name: "id", Number: 1, TypeName: "string"}}},
					{Name: "GetResp", Fields: []TestField{{Name: "name", Number: 1, TypeName: "string"}}},
				},
				Services: []TestService{
					{
						Name: "UserService",
						Methods: []TestMethod{
							{
								Name:       "GetUser",
								InputType:  "test.v1.GetReq",
								OutputType: "test.v1.GetResp",
							},
							{
								Name:            "StreamUsers",
								InputType:       "test.v1.GetReq",
								OutputType:      "test.v1.GetResp",
								ServerStreaming:  true,
							},
						},
					},
				},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	file := plugin.Files[0]

	if len(file.Services) != 1 {
		t.Fatalf("Expected 1 service, got %d", len(file.Services))
	}

	svc := file.Services[0]
	if string(svc.Desc.Name()) != "UserService" {
		t.Errorf("Expected service name 'UserService', got '%s'", svc.Desc.Name())
	}
	if len(svc.Methods) != 2 {
		t.Fatalf("Expected 2 methods, got %d", len(svc.Methods))
	}

	getUser := svc.Methods[0]
	if string(getUser.Desc.Name()) != "GetUser" {
		t.Errorf("Expected method 'GetUser', got '%s'", getUser.Desc.Name())
	}
	if getUser.Desc.IsStreamingClient() || getUser.Desc.IsStreamingServer() {
		t.Error("GetUser should not be streaming")
	}

	streamUsers := svc.Methods[1]
	if !streamUsers.Desc.IsStreamingServer() {
		t.Error("StreamUsers should be server-streaming")
	}
	if streamUsers.Desc.IsStreamingClient() {
		t.Error("StreamUsers should not be client-streaming")
	}
}

func TestAllScalarTypes(t *testing.T) {
	// Verify all scalar type mappings work.
	scalars := []string{"string", "int32", "sint32", "sfixed32", "int64", "sint64", "sfixed64",
		"uint32", "fixed32", "uint64", "fixed64", "bool", "float", "double", "bytes"}

	var testFields []TestField
	for i, s := range scalars {
		testFields = append(testFields, TestField{Name: s + "_field", Number: int32(i + 1), TypeName: s})
	}

	protoSet := &TestProtoSet{
		Files: []TestFile{
			{
				Name:     "test.proto",
				Pkg:      "test.v1",
				Messages: []TestMessage{{Name: "AllScalars", Fields: testFields}},
			},
		},
	}

	plugin := CreateTestPlugin(t, protoSet)
	msg := plugin.Files[0].Messages[0]

	if len(msg.Fields) != len(scalars) {
		t.Fatalf("Expected %d fields, got %d", len(scalars), len(msg.Fields))
	}
}
