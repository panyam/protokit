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

package wellknown

import (
	"maps"

	"google.golang.org/protobuf/compiler/protogen"
)

// TypeMapping represents a mapping from a proto well-known type to a target language type.
type TypeMapping struct {
	ProtoFullName string // e.g., "google.protobuf.Timestamp"
	TargetType    string // e.g., "time.Time" or "Timestamp" — plugin decides
	ImportPath    string // e.g., "time" or "@bufbuild/protobuf/wkt"
	IsNative      bool   // Maps to a built-in type in target language
}

// Registry is a language-agnostic well-known type registry.
// Plugins register their own target-language mappings.
type Registry struct {
	mappings map[string]TypeMapping
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{
		mappings: make(map[string]TypeMapping),
	}
}

// Register adds a type mapping.
func (r *Registry) Register(protoFullName, targetType, importPath string, isNative bool) {
	r.mappings[protoFullName] = TypeMapping{
		ProtoFullName: protoFullName,
		TargetType:    targetType,
		ImportPath:    importPath,
		IsNative:      isNative,
	}
}

// Get returns the mapping for a proto type, if registered.
func (r *Registry) Get(protoFullName string) (TypeMapping, bool) {
	m, ok := r.mappings[protoFullName]
	return m, ok
}

// GetByMessage returns the mapping for a protogen message's fully qualified name.
func (r *Registry) GetByMessage(msg *protogen.Message) (TypeMapping, bool) {
	return r.Get(string(msg.Desc.FullName()))
}

// IsWellKnown checks if a proto type is registered.
func (r *Registry) IsWellKnown(protoFullName string) bool {
	_, ok := r.mappings[protoFullName]
	return ok
}

// AllMappings returns a copy of all registered mappings.
func (r *Registry) AllMappings() map[string]TypeMapping {
	result := make(map[string]TypeMapping, len(r.mappings))
	maps.Copy(result, r.mappings)
	return result
}

// WellKnownProtoTypes returns the list of standard google.protobuf.* type names.
// This is a convenience — plugins can use this to initialize their registries.
func WellKnownProtoTypes() []string {
	return []string{
		"google.protobuf.Timestamp",
		"google.protobuf.Duration",
		"google.protobuf.Any",
		"google.protobuf.Empty",
		"google.protobuf.Struct",
		"google.protobuf.Value",
		"google.protobuf.ListValue",
		"google.protobuf.NullValue",
		"google.protobuf.FieldMask",
		"google.protobuf.DoubleValue",
		"google.protobuf.FloatValue",
		"google.protobuf.Int64Value",
		"google.protobuf.UInt64Value",
		"google.protobuf.Int32Value",
		"google.protobuf.UInt32Value",
		"google.protobuf.BoolValue",
		"google.protobuf.StringValue",
		"google.protobuf.BytesValue",
		"google.protobuf.Type",
		"google.protobuf.Field",
		"google.protobuf.Enum",
		"google.protobuf.EnumValue",
		"google.protobuf.Option",
		"google.protobuf.SourceContext",
	}
}
