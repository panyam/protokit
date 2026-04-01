// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
