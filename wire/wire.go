// Package wire provides low-level helpers for decoding proto wire-format bytes.
//
// These are useful when extracting proto extension fields from descriptor options
// without a full proto unmarshal (e.g., reading custom method/service annotations
// from protoc-gen plugins).
//
// All Decode* functions scan raw wire bytes for a specific field number and
// return the decoded value, or the zero value if the field is not present.
package wire

import (
	"math"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ExtractExtension marshals a proto message and extracts the raw bytes of
// a specific extension field. Returns nil if the message is nil or the
// field is not present.
func ExtractExtension(msg protoreflect.ProtoMessage, field protowire.Number) []byte {
	if msg == nil {
		return nil
	}
	raw, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}
	return DecodeBytes(raw, field)
}

// DecodeString scans raw proto bytes for a length-delimited field and returns
// it as a string. Returns "" if not found.
func DecodeString(raw []byte, fieldNum protowire.Number) string {
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return ""
		}
		raw = raw[n:]
		n = skipOrExtract(raw, typ)
		if n < 0 {
			return ""
		}
		if num == fieldNum && typ == protowire.BytesType {
			v, _ := protowire.ConsumeBytes(raw)
			return string(v)
		}
		raw = raw[n:]
	}
	return ""
}

// DecodeStringList collects all string values for a repeated length-delimited field.
func DecodeStringList(raw []byte, fieldNum protowire.Number) []string {
	var result []string
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return result
		}
		raw = raw[n:]
		n = skipOrExtract(raw, typ)
		if n < 0 {
			return result
		}
		if num == fieldNum && typ == protowire.BytesType {
			v, _ := protowire.ConsumeBytes(raw)
			result = append(result, string(v))
		}
		raw = raw[n:]
	}
	return result
}

// DecodeBool scans raw proto bytes for a varint field and returns it as a bool.
// Returns false if not found.
func DecodeBool(raw []byte, fieldNum protowire.Number) bool {
	v, ok := decodeVarint(raw, fieldNum)
	if !ok {
		return false
	}
	return v != 0
}

// DecodeInt32 scans raw proto bytes for a varint field and returns it as int32.
// Returns 0 if not found.
func DecodeInt32(raw []byte, fieldNum protowire.Number) int32 {
	v, ok := decodeVarint(raw, fieldNum)
	if !ok {
		return 0
	}
	return int32(v)
}

// DecodeFloat scans raw proto bytes for a fixed32 field and returns it as float32.
// Returns 0 if not found.
func DecodeFloat(raw []byte, fieldNum protowire.Number) float32 {
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return 0
		}
		raw = raw[n:]
		n = skipOrExtract(raw, typ)
		if n < 0 {
			return 0
		}
		if num == fieldNum && typ == protowire.Fixed32Type {
			v, _ := protowire.ConsumeFixed32(raw)
			return math.Float32frombits(v)
		}
		raw = raw[n:]
	}
	return 0
}

// DecodeBytes scans raw proto bytes for a length-delimited field and returns
// the raw bytes. Returns nil if not found. Use this for embedded messages.
func DecodeBytes(raw []byte, fieldNum protowire.Number) []byte {
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return nil
		}
		raw = raw[n:]
		n = skipOrExtract(raw, typ)
		if n < 0 {
			return nil
		}
		if num == fieldNum && typ == protowire.BytesType {
			v, _ := protowire.ConsumeBytes(raw)
			return v
		}
		raw = raw[n:]
	}
	return nil
}

// decodeVarint is a shared helper that scans for a varint field.
func decodeVarint(raw []byte, fieldNum protowire.Number) (uint64, bool) {
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return 0, false
		}
		raw = raw[n:]
		n = skipOrExtract(raw, typ)
		if n < 0 {
			return 0, false
		}
		if num == fieldNum && typ == protowire.VarintType {
			v, _ := protowire.ConsumeVarint(raw)
			return v, true
		}
		raw = raw[n:]
	}
	return 0, false
}

// skipOrExtract returns the number of bytes to skip for the given wire type.
func skipOrExtract(raw []byte, typ protowire.Type) int {
	switch typ {
	case protowire.VarintType:
		_, n := protowire.ConsumeVarint(raw)
		return n
	case protowire.Fixed32Type:
		_, n := protowire.ConsumeFixed32(raw)
		return n
	case protowire.Fixed64Type:
		_, n := protowire.ConsumeFixed64(raw)
		return n
	case protowire.BytesType:
		_, n := protowire.ConsumeBytes(raw)
		return n
	default:
		return -1
	}
}
