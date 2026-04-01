// Copyright 2024 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fields

import "google.golang.org/protobuf/compiler/protogen"

// GetFieldKind returns the proto field kind as a string (e.g., "string", "int32", "message", "enum").
func GetFieldKind(field *protogen.Field) string {
	return field.Desc.Kind().String()
}

// IsMapField checks if a field is a proto map type.
func IsMapField(field *protogen.Field) bool {
	return field.Message != nil && field.Message.Desc.IsMapEntry()
}

// GetMapKeyValueFields extracts the key and value fields from a map field.
// Returns nil, nil if not a map field.
func GetMapKeyValueFields(field *protogen.Field) (keyField, valueField *protogen.Field) {
	if !IsMapField(field) {
		return nil, nil
	}
	for _, f := range field.Message.Fields {
		switch f.Desc.Number() {
		case 1:
			keyField = f
		case 2:
			valueField = f
		}
	}
	return keyField, valueField
}

// IsRepeated checks if a field is a repeated (list) field.
func IsRepeated(field *protogen.Field) bool {
	return field.Desc.IsList()
}

// IsOptional checks if a field has the proto3 optional keyword.
func IsOptional(field *protogen.Field) bool {
	return field.Desc.HasOptionalKeyword()
}

// IsNumericKind checks if a proto kind string represents a numeric type.
func IsNumericKind(kind string) bool {
	switch kind {
	case "int32", "int64", "uint32", "uint64",
		"sint32", "sint64", "fixed32", "fixed64",
		"sfixed32", "sfixed64", "float", "double":
		return true
	}
	return false
}

// NormalizeNumericKind maps aliased numeric kinds to their canonical Go types.
// e.g., "sint32" -> "int32", "fixed64" -> "uint64", "double" -> "float64"
func NormalizeNumericKind(kind string) string {
	switch kind {
	case "sint32":
		return "int32"
	case "sint64":
		return "int64"
	case "fixed32":
		return "uint32"
	case "fixed64":
		return "uint64"
	case "sfixed32":
		return "int32"
	case "sfixed64":
		return "int64"
	case "float":
		return "float32"
	case "double":
		return "float64"
	default:
		return kind
	}
}
