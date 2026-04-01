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

package messages

import (
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ExtractPackageName extracts the package from a fully qualified type.
// "library.v1.Book" -> "library.v1"
func ExtractPackageName(fullType string) string {
	lastDot := strings.LastIndex(fullType, ".")
	if lastDot < 0 {
		return ""
	}
	return fullType[:lastDot]
}

// ExtractMessageName extracts the message name from a fully qualified type.
// "library.v1.Book" -> "Book"
func ExtractMessageName(fullType string) string {
	lastDot := strings.LastIndex(fullType, ".")
	if lastDot < 0 {
		return fullType
	}
	return fullType[lastDot+1:]
}

// GetFullyQualifiedType returns the fully qualified type for a message field.
// Returns "" if field.Message is nil.
func GetFullyQualifiedType(field *protogen.Field) string {
	if field.Message == nil {
		return ""
	}
	pkg := string(field.Message.Desc.ParentFile().Package())
	name := string(field.Message.Desc.Name())
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

// IsNestedMessage checks if a message is nested inside another message.
func IsNestedMessage(message *protogen.Message) bool {
	_, isNestedInMessage := message.Desc.Parent().(protoreflect.MessageDescriptor)
	return isNestedInMessage
}

// GetOneofGroups returns the names of all oneof groups in a message.
func GetOneofGroups(message *protogen.Message) []string {
	var names []string
	for _, oneof := range message.Oneofs {
		names = append(names, string(oneof.Desc.Name()))
	}
	return names
}

// GetBaseFileName extracts the base filename from a proto file path.
// "path/to/library.proto" -> "library"
func GetBaseFileName(protoFile string) string {
	base := filepath.Base(protoFile)
	return strings.TrimSuffix(base, ".proto")
}

// BuildMessageIndex builds a map from fully qualified message name to *protogen.Message
// across all files in the plugin (including non-generated imports).
func BuildMessageIndex(gen *protogen.Plugin) map[string]*protogen.Message {
	index := make(map[string]*protogen.Message)
	for _, file := range gen.Files {
		indexMessages(file.Messages, index)
	}
	return index
}

// indexMessages recursively indexes messages and their nested messages.
func indexMessages(messages []*protogen.Message, index map[string]*protogen.Message) {
	for _, msg := range messages {
		index[string(msg.Desc.FullName())] = msg
		indexMessages(msg.Messages, index)
	}
}
