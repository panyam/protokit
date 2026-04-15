package wire

import (
	"math"
	"testing"

	"google.golang.org/protobuf/encoding/protowire"
)

// buildMessage constructs raw proto bytes from a set of field values.
func buildMessage(fields map[protowire.Number]any) []byte {
	var buf []byte
	for num, val := range fields {
		switch v := val.(type) {
		case string:
			buf = protowire.AppendTag(buf, num, protowire.BytesType)
			buf = protowire.AppendString(buf, v)
		case bool:
			buf = protowire.AppendTag(buf, num, protowire.VarintType)
			if v {
				buf = protowire.AppendVarint(buf, 1)
			} else {
				buf = protowire.AppendVarint(buf, 0)
			}
		case int32:
			buf = protowire.AppendTag(buf, num, protowire.VarintType)
			buf = protowire.AppendVarint(buf, uint64(v))
		case float32:
			buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
			buf = protowire.AppendFixed32(buf, math.Float32bits(v))
		case []byte:
			buf = protowire.AppendTag(buf, num, protowire.BytesType)
			buf = protowire.AppendBytes(buf, v)
		}
	}
	return buf
}

func TestDecodeString(t *testing.T) {
	raw := buildMessage(map[protowire.Number]any{
		1: "hello",
		2: "world",
	})

	if got := DecodeString(raw, 1); got != "hello" {
		t.Errorf("DecodeString(1) = %q, want %q", got, "hello")
	}
	if got := DecodeString(raw, 2); got != "world" {
		t.Errorf("DecodeString(2) = %q, want %q", got, "world")
	}
	if got := DecodeString(raw, 3); got != "" {
		t.Errorf("DecodeString(3) = %q, want empty", got)
	}
	if got := DecodeString(nil, 1); got != "" {
		t.Errorf("DecodeString(nil) = %q, want empty", got)
	}
}

func TestDecodeBool(t *testing.T) {
	raw := buildMessage(map[protowire.Number]any{
		1: true,
		2: false,
	})

	if got := DecodeBool(raw, 1); !got {
		t.Error("DecodeBool(1) = false, want true")
	}
	if got := DecodeBool(raw, 2); got {
		t.Error("DecodeBool(2) = true, want false")
	}
	if got := DecodeBool(raw, 3); got {
		t.Error("DecodeBool(3) = true, want false")
	}
}

func TestDecodeInt32(t *testing.T) {
	raw := buildMessage(map[protowire.Number]any{
		1: int32(42),
		2: int32(0),
	})

	if got := DecodeInt32(raw, 1); got != 42 {
		t.Errorf("DecodeInt32(1) = %d, want 42", got)
	}
	if got := DecodeInt32(raw, 2); got != 0 {
		t.Errorf("DecodeInt32(2) = %d, want 0", got)
	}
	if got := DecodeInt32(raw, 3); got != 0 {
		t.Errorf("DecodeInt32(3) = %d, want 0", got)
	}
}

func TestDecodeFloat(t *testing.T) {
	raw := buildMessage(map[protowire.Number]any{
		1: float32(0.8),
		2: float32(0.0),
	})

	if got := DecodeFloat(raw, 1); got != 0.8 {
		t.Errorf("DecodeFloat(1) = %f, want 0.8", got)
	}
	if got := DecodeFloat(raw, 2); got != 0.0 {
		t.Errorf("DecodeFloat(2) = %f, want 0.0", got)
	}
	if got := DecodeFloat(raw, 3); got != 0.0 {
		t.Errorf("DecodeFloat(3) = %f, want 0.0", got)
	}
}

func TestDecodeBytes(t *testing.T) {
	inner := buildMessage(map[protowire.Number]any{1: "nested"})
	raw := buildMessage(map[protowire.Number]any{
		1: "text",
		2: inner, // embedded message as []byte
	})

	got := DecodeBytes(raw, 2)
	if got == nil {
		t.Fatal("DecodeBytes(2) = nil, want nested message bytes")
	}
	if s := DecodeString(got, 1); s != "nested" {
		t.Errorf("nested DecodeString(1) = %q, want %q", s, "nested")
	}
	if got := DecodeBytes(raw, 3); got != nil {
		t.Error("DecodeBytes(3) should be nil for missing field")
	}
}

func TestDecodeStringList(t *testing.T) {
	// Build repeated field manually (same field number, multiple values).
	var raw []byte
	raw = protowire.AppendTag(raw, 1, protowire.BytesType)
	raw = protowire.AppendString(raw, "alpha")
	raw = protowire.AppendTag(raw, 1, protowire.BytesType)
	raw = protowire.AppendString(raw, "beta")
	raw = protowire.AppendTag(raw, 2, protowire.BytesType)
	raw = protowire.AppendString(raw, "other")

	got := DecodeStringList(raw, 1)
	if len(got) != 2 || got[0] != "alpha" || got[1] != "beta" {
		t.Errorf("DecodeStringList(1) = %v, want [alpha beta]", got)
	}
	got = DecodeStringList(raw, 2)
	if len(got) != 1 || got[0] != "other" {
		t.Errorf("DecodeStringList(2) = %v, want [other]", got)
	}
	got = DecodeStringList(raw, 3)
	if len(got) != 0 {
		t.Errorf("DecodeStringList(3) = %v, want empty", got)
	}
}
