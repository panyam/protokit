// Package testutil provides helpers for building in-memory proto descriptors
// for unit testing protoc plugins without requiring actual .proto files.
package testutil

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// TestProtoSet represents a collection of proto files for testing.
type TestProtoSet struct {
	Files []TestFile
}

// TestFile represents a single proto file with messages.
type TestFile struct {
	Name     string
	Pkg      string
	Messages []TestMessage
	Enums    []TestEnum
	Services []TestService
}

// TestMessage represents a proto message with optional plugin-specific options.
type TestMessage struct {
	Name    string
	Fields  []TestField
	Oneofs  []string                      // oneof group names
	Options *descriptorpb.MessageOptions  // Generic — plugins set their own extensions
}

// TestField represents a proto field with optional plugin-specific options.
type TestField struct {
	Name       string
	Number     int32
	TypeName   string                          // "string", "int32", "int64", "uint32", "uint64", "bool", "float", "double", "bytes", or fully qualified message name
	Repeated   bool
	Optional   bool                            // proto3 optional keyword
	IsMap      bool
	MapKeyType string                          // For map fields: "int32", "string", etc.
	OneofIndex int                             // -1 or unset = not in a oneof, otherwise index into Message.Oneofs
	EnumType   string                          // fully qualified enum type name (e.g. "test.Status")
	Options    *descriptorpb.FieldOptions      // Generic — plugins set their own extensions
}

// TestEnum represents a proto enum.
type TestEnum struct {
	Name   string
	Values []TestEnumValue
}

// TestEnumValue represents a proto enum value.
type TestEnumValue struct {
	Name   string
	Number int32
}

// TestService represents a proto service with methods.
type TestService struct {
	Name    string
	Methods []TestMethod
}

// TestMethod represents a service RPC method.
type TestMethod struct {
	Name            string
	InputType       string // fully qualified message name (e.g. "pkg.v1.GetUserRequest")
	OutputType      string // fully qualified message name
	ClientStreaming bool
	ServerStreaming bool
}

// CreateTestPlugin creates a protogen.Plugin from a test proto set.
func CreateTestPlugin(t *testing.T, protoSet *TestProtoSet) *protogen.Plugin {
	t.Helper()

	req := BuildCodeGeneratorRequest(t, protoSet)
	opts := protogen.Options{}
	plugin, err := opts.New(req)
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	return plugin
}

// BuildCodeGeneratorRequest creates a CodeGeneratorRequest from a test proto set.
func BuildCodeGeneratorRequest(t *testing.T, protoSet *TestProtoSet) *pluginpb.CodeGeneratorRequest {
	t.Helper()

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{},
		ProtoFile:      []*descriptorpb.FileDescriptorProto{},
	}

	for _, file := range protoSet.Files {
		fileDesc := BuildFileDescriptor(t, file)
		req.ProtoFile = append(req.ProtoFile, fileDesc)
		req.FileToGenerate = append(req.FileToGenerate, file.Name)
	}

	return req
}

// BuildFileDescriptor creates a FileDescriptorProto from a test file.
func BuildFileDescriptor(t *testing.T, file TestFile) *descriptorpb.FileDescriptorProto {
	t.Helper()

	goPackage := "github.com/test/gen/go/" + strings.ReplaceAll(file.Pkg, ".", "/")

	fileDesc := &descriptorpb.FileDescriptorProto{
		Name:    proto.String(file.Name),
		Package: proto.String(file.Pkg),
		Syntax:  proto.String("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String(goPackage),
		},
	}

	for _, e := range file.Enums {
		fileDesc.EnumType = append(fileDesc.EnumType, buildEnumDescriptor(e))
	}

	for _, msg := range file.Messages {
		msgDesc := BuildMessageDescriptorWithPackage(t, msg, file.Pkg)
		fileDesc.MessageType = append(fileDesc.MessageType, msgDesc)
	}

	for _, svc := range file.Services {
		fileDesc.Service = append(fileDesc.Service, buildServiceDescriptor(svc))
	}

	return fileDesc
}

// BuildMessageDescriptor creates a DescriptorProto from a test message (no package context).
func BuildMessageDescriptor(t *testing.T, msg TestMessage) *descriptorpb.DescriptorProto {
	return BuildMessageDescriptorWithPackage(t, msg, "")
}

// BuildMessageDescriptorWithPackage creates a DescriptorProto from a test message with package context.
func BuildMessageDescriptorWithPackage(t *testing.T, msg TestMessage, pkg string) *descriptorpb.DescriptorProto {
	t.Helper()

	msgDesc := &descriptorpb.DescriptorProto{
		Name:    proto.String(msg.Name),
		Options: msg.Options,
	}

	// Add declared oneof groups.
	for _, name := range msg.Oneofs {
		msgDesc.OneofDecl = append(msgDesc.OneofDecl, &descriptorpb.OneofDescriptorProto{
			Name: proto.String(name),
		})
	}

	hasOneofs := len(msg.Oneofs) > 0
	syntheticOneofIndex := int32(len(msg.Oneofs))

	// Add fields
	for _, field := range msg.Fields {
		if field.IsMap {
			addMapField(msgDesc, field, msg.Name, pkg)
		} else {
			fieldDesc := buildFieldDescriptor(field, hasOneofs, syntheticOneofIndex)
			msgDesc.Field = append(msgDesc.Field, fieldDesc)

			// proto3 optional uses a synthetic oneof.
			if field.Optional {
				syntheticName := "_" + field.Name
				msgDesc.OneofDecl = append(msgDesc.OneofDecl, &descriptorpb.OneofDescriptorProto{
					Name: proto.String(syntheticName),
				})
				syntheticOneofIndex++
			}
		}
	}

	return msgDesc
}

func buildFieldDescriptor(field TestField, hasOneofs bool, syntheticOneofBase int32) *descriptorpb.FieldDescriptorProto {
	fieldDesc := &descriptorpb.FieldDescriptorProto{
		Name:    proto.String(field.Name),
		Number:  proto.Int32(field.Number),
		Options: field.Options,
		Label:   descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
	}

	if field.EnumType != "" {
		fieldDesc.Type = descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum()
		fieldDesc.TypeName = proto.String("." + field.EnumType)
	} else {
		fieldDesc.Type = GetFieldType(field.TypeName)
		if typeName := GetTypeName(field.TypeName); typeName != nil {
			fieldDesc.TypeName = typeName
		}
	}

	// Set label for repeated fields
	if field.Repeated {
		fieldDesc.Label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()
	}

	// Handle oneof membership. OneofIndex >= 0 means the field belongs to
	// the oneof at that index. Negative values (or zero when no oneofs declared)
	// mean the field is not in any oneof.
	if !field.Optional && field.OneofIndex >= 0 && hasOneofs {
		idx := int32(field.OneofIndex)
		fieldDesc.OneofIndex = &idx
	}

	// Handle proto3 optional (synthetic oneof).
	if field.Optional {
		fieldDesc.Proto3Optional = proto.Bool(true)
		fieldDesc.OneofIndex = &syntheticOneofBase
	}

	return fieldDesc
}

func addMapField(msgDesc *descriptorpb.DescriptorProto, field TestField, msgName, pkg string) {
	// Map fields require a nested entry message.
	fieldName := field.Name
	if len(fieldName) > 0 {
		fieldName = strings.ToUpper(fieldName[:1]) + fieldName[1:]
	}
	entryMsgName := fieldName + "Entry"

	valFd := &descriptorpb.FieldDescriptorProto{
		Name:   proto.String("value"),
		Number: proto.Int32(2),
		Type:   GetFieldType(field.TypeName),
		Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
	}
	if typeName := GetTypeName(field.TypeName); typeName != nil {
		valFd.TypeName = typeName
	}

	entryMsg := &descriptorpb.DescriptorProto{
		Name: proto.String(entryMsgName),
		Options: &descriptorpb.MessageOptions{
			MapEntry: proto.Bool(true),
		},
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("key"),
				Number: proto.Int32(1),
				Type:   GetFieldType(field.MapKeyType),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
			},
			valFd,
		},
	}
	msgDesc.NestedType = append(msgDesc.NestedType, entryMsg)

	// Add the map field itself.
	fullEntryName := "." + msgName + "." + entryMsgName
	if pkg != "" {
		fullEntryName = "." + pkg + "." + msgName + "." + entryMsgName
	}
	fieldDesc := &descriptorpb.FieldDescriptorProto{
		Name:     proto.String(field.Name),
		Number:   proto.Int32(field.Number),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		TypeName: proto.String(fullEntryName),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
		Options:  field.Options,
	}
	msgDesc.Field = append(msgDesc.Field, fieldDesc)
}

func buildEnumDescriptor(e TestEnum) *descriptorpb.EnumDescriptorProto {
	ed := &descriptorpb.EnumDescriptorProto{Name: proto.String(e.Name)}
	for _, v := range e.Values {
		ed.Value = append(ed.Value, &descriptorpb.EnumValueDescriptorProto{
			Name:   proto.String(v.Name),
			Number: proto.Int32(v.Number),
		})
	}
	return ed
}

func buildServiceDescriptor(svc TestService) *descriptorpb.ServiceDescriptorProto {
	sd := &descriptorpb.ServiceDescriptorProto{Name: proto.String(svc.Name)}
	for _, m := range svc.Methods {
		sd.Method = append(sd.Method, &descriptorpb.MethodDescriptorProto{
			Name:            proto.String(m.Name),
			InputType:       proto.String("." + m.InputType),
			OutputType:      proto.String("." + m.OutputType),
			ClientStreaming: proto.Bool(m.ClientStreaming),
			ServerStreaming: proto.Bool(m.ServerStreaming),
		})
	}
	return sd
}

// GetFieldType returns the proto field type enum for a type name string.
func GetFieldType(typeName string) *descriptorpb.FieldDescriptorProto_Type {
	switch typeName {
	case "string":
		return descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()
	case "int32":
		return descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum()
	case "sint32":
		return descriptorpb.FieldDescriptorProto_TYPE_SINT32.Enum()
	case "sfixed32":
		return descriptorpb.FieldDescriptorProto_TYPE_SFIXED32.Enum()
	case "int64":
		return descriptorpb.FieldDescriptorProto_TYPE_INT64.Enum()
	case "sint64":
		return descriptorpb.FieldDescriptorProto_TYPE_SINT64.Enum()
	case "sfixed64":
		return descriptorpb.FieldDescriptorProto_TYPE_SFIXED64.Enum()
	case "uint32":
		return descriptorpb.FieldDescriptorProto_TYPE_UINT32.Enum()
	case "fixed32":
		return descriptorpb.FieldDescriptorProto_TYPE_FIXED32.Enum()
	case "uint64":
		return descriptorpb.FieldDescriptorProto_TYPE_UINT64.Enum()
	case "fixed64":
		return descriptorpb.FieldDescriptorProto_TYPE_FIXED64.Enum()
	case "bool":
		return descriptorpb.FieldDescriptorProto_TYPE_BOOL.Enum()
	case "float":
		return descriptorpb.FieldDescriptorProto_TYPE_FLOAT.Enum()
	case "double":
		return descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum()
	case "bytes":
		return descriptorpb.FieldDescriptorProto_TYPE_BYTES.Enum()
	default:
		return descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum()
	}
}

// GetTypeName returns the fully qualified type name for message types, nil for scalars.
func GetTypeName(typeName string) *string {
	switch typeName {
	case "string", "int32", "sint32", "sfixed32", "int64", "sint64", "sfixed64",
		"uint32", "fixed32", "uint64", "fixed64", "bool", "float", "double", "bytes":
		return nil
	default:
		return proto.String("." + typeName)
	}
}
