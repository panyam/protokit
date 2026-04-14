package fields

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// PopulateFieldFromPath sets a field on a proto message from a string value.
// The fieldPath supports dot-separated paths for nested fields (e.g., "pos.q").
// Type coercion from string to the target field type is handled automatically
// for all scalar types and enums.
func PopulateFieldFromPath(msg proto.Message, fieldPath string, value string) error {
	parts := strings.Split(fieldPath, ".")
	return populatePath(msg.ProtoReflect(), parts, value)
}

// PopulateFromMap sets multiple fields on a proto message from a string map.
// Each key is a field path (dot-separated for nested fields) and each value
// is coerced to the target field type.
func PopulateFromMap(msg proto.Message, params map[string]string) error {
	for key, value := range params {
		if err := PopulateFieldFromPath(msg, key, value); err != nil {
			return err
		}
	}
	return nil
}

func populatePath(m protoreflect.Message, path []string, value string) error {
	md := m.Descriptor()
	fd := md.Fields().ByName(protoreflect.Name(path[0]))
	if fd == nil {
		return fmt.Errorf("field %q not found on %s", path[0], md.FullName())
	}

	if len(path) > 1 {
		// Nested path — descend into sub-message.
		if fd.Kind() != protoreflect.MessageKind {
			return fmt.Errorf("field %q on %s is not a message (got %s), cannot traverse path", path[0], md.FullName(), fd.Kind())
		}
		sub := m.Mutable(fd).Message()
		return populatePath(sub, path[1:], value)
	}

	// Leaf field — parse and set the value.
	v, err := parseFieldValue(fd, value)
	if err != nil {
		return fmt.Errorf("field %q on %s: %w", path[0], md.FullName(), err)
	}
	m.Set(fd, v)
	return nil
}

func parseFieldValue(fd protoreflect.FieldDescriptor, s string) (protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(s), nil

	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(s)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as bool: %w", s, err)
		}
		return protoreflect.ValueOfBool(v), nil

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as int32: %w", s, err)
		}
		return protoreflect.ValueOfInt32(int32(v)), nil

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as int64: %w", s, err)
		}
		return protoreflect.ValueOfInt64(v), nil

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as uint32: %w", s, err)
		}
		return protoreflect.ValueOfUint32(uint32(v)), nil

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as uint64: %w", s, err)
		}
		return protoreflect.ValueOfUint64(v), nil

	case protoreflect.FloatKind:
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as float: %w", s, err)
		}
		return protoreflect.ValueOfFloat32(float32(v)), nil

	case protoreflect.DoubleKind:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as double: %w", s, err)
		}
		return protoreflect.ValueOfFloat64(v), nil

	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes([]byte(s)), nil

	case protoreflect.EnumKind:
		// Try name first, then numeric value.
		ev := fd.Enum().Values().ByName(protoreflect.Name(s))
		if ev != nil {
			return protoreflect.ValueOfEnum(ev.Number()), nil
		}
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("cannot parse %q as enum %s: not a valid name or number", s, fd.Enum().FullName())
		}
		return protoreflect.ValueOfEnum(protoreflect.EnumNumber(n)), nil

	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported field kind %s for string coercion", fd.Kind())
	}
}
