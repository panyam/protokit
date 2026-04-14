package fields_test

import (
	"testing"

	"github.com/panyam/protokit/fields"
	"github.com/panyam/protokit/testutil"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// buildDynMsg creates a dynamic proto message from a test message definition.
// Returns the message descriptor and a fresh instance.
func buildDynMsg(t *testing.T, msgs []testutil.TestMessage) (protoreflect.MessageDescriptor, *dynamicpb.Message) {
	t.Helper()
	plugin := testutil.CreateTestPlugin(t, &testutil.TestProtoSet{
		Files: []testutil.TestFile{{
			Name: "test.proto",
			Pkg:  "test.v1",
			Messages: msgs,
		}},
	})
	for _, f := range plugin.Files {
		for _, m := range f.Messages {
			if string(m.Desc.Name()) == msgs[0].Name {
				return m.Desc, dynamicpb.NewMessage(m.Desc)
			}
		}
	}
	t.Fatal("message not found")
	return nil, nil
}

// TestPopulateStringField verifies setting a string field from a string value.
func TestPopulateStringField(t *testing.T) {
	md, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "user_id", Number: 1, TypeName: "string"},
		},
	}})

	err := fields.PopulateFieldFromPath(msg, "user_id", "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fd := md.Fields().ByName("user_id")
	got := msg.Get(fd).String()
	if got != "abc-123" {
		t.Errorf("user_id = %q, want %q", got, "abc-123")
	}
}

// TestPopulateInt32Field verifies string-to-int32 coercion.
func TestPopulateInt32Field(t *testing.T) {
	md, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "count", Number: 1, TypeName: "int32"},
		},
	}})

	err := fields.PopulateFieldFromPath(msg, "count", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fd := md.Fields().ByName("count")
	got := msg.Get(fd).Int()
	if got != 42 {
		t.Errorf("count = %d, want 42", got)
	}
}

// TestPopulateInt64Field verifies string-to-int64 coercion.
func TestPopulateInt64Field(t *testing.T) {
	md, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "big_id", Number: 1, TypeName: "int64"},
		},
	}})

	err := fields.PopulateFieldFromPath(msg, "big_id", "9999999999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fd := md.Fields().ByName("big_id")
	got := msg.Get(fd).Int()
	if got != 9999999999 {
		t.Errorf("big_id = %d, want 9999999999", got)
	}
}

// TestPopulateBoolField verifies string-to-bool coercion.
func TestPopulateBoolField(t *testing.T) {
	md, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "active", Number: 1, TypeName: "bool"},
		},
	}})

	err := fields.PopulateFieldFromPath(msg, "active", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fd := md.Fields().ByName("active")
	got := msg.Get(fd).Bool()
	if !got {
		t.Error("active = false, want true")
	}
}

// TestPopulateFloatField verifies string-to-float coercion.
func TestPopulateFloatField(t *testing.T) {
	md, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "score", Number: 1, TypeName: "float"},
		},
	}})

	err := fields.PopulateFieldFromPath(msg, "score", "3.14")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fd := md.Fields().ByName("score")
	got := msg.Get(fd).Float()
	if got < 3.13 || got > 3.15 {
		t.Errorf("score = %f, want ~3.14", got)
	}
}

// TestPopulateNestedDotPath verifies dot-separated paths for nested fields
// (e.g., "pos.q" sets the q field inside a pos sub-message).
func TestPopulateNestedDotPath(t *testing.T) {
	// Build Position and Req messages.
	plugin := testutil.CreateTestPlugin(t, &testutil.TestProtoSet{
		Files: []testutil.TestFile{{
			Name: "test.proto",
			Pkg:  "test.v1",
			Messages: []testutil.TestMessage{
				{Name: "Position", Fields: []testutil.TestField{
					{Name: "q", Number: 1, TypeName: "int32"},
					{Name: "r", Number: 2, TypeName: "int32"},
				}},
				{Name: "Req", Fields: []testutil.TestField{
					{Name: "game_id", Number: 1, TypeName: "string"},
					{Name: "pos", Number: 2, TypeName: "test.v1.Position"},
				}},
			},
		}},
	})

	var reqDesc protoreflect.MessageDescriptor
	for _, f := range plugin.Files {
		for _, m := range f.Messages {
			if string(m.Desc.Name()) == "Req" {
				reqDesc = m.Desc
			}
		}
	}
	if reqDesc == nil {
		t.Fatal("Req message not found")
	}

	msg := dynamicpb.NewMessage(reqDesc)

	// Set flat field.
	if err := fields.PopulateFieldFromPath(msg, "game_id", "game-42"); err != nil {
		t.Fatalf("game_id: %v", err)
	}
	// Set nested fields.
	if err := fields.PopulateFieldFromPath(msg, "pos.q", "3"); err != nil {
		t.Fatalf("pos.q: %v", err)
	}
	if err := fields.PopulateFieldFromPath(msg, "pos.r", "5"); err != nil {
		t.Fatalf("pos.r: %v", err)
	}

	// Verify.
	gameID := msg.Get(reqDesc.Fields().ByName("game_id")).String()
	if gameID != "game-42" {
		t.Errorf("game_id = %q, want %q", gameID, "game-42")
	}

	posFd := reqDesc.Fields().ByName("pos")
	posMsg := msg.Get(posFd).Message()
	posDesc := posFd.Message()
	q := posMsg.Get(posDesc.Fields().ByName("q")).Int()
	r := posMsg.Get(posDesc.Fields().ByName("r")).Int()
	if q != 3 {
		t.Errorf("pos.q = %d, want 3", q)
	}
	if r != 5 {
		t.Errorf("pos.r = %d, want 5", r)
	}
}

// TestPopulateFromMap verifies setting multiple fields at once from a map.
func TestPopulateFromMap(t *testing.T) {
	md, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "design_id", Number: 1, TypeName: "string"},
			{Name: "section_id", Number: 2, TypeName: "string"},
			{Name: "name", Number: 3, TypeName: "string"},
		},
	}})

	err := fields.PopulateFromMap(msg, map[string]string{
		"design_id":  "d1",
		"section_id": "s2",
		"name":       "main",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, tc := range []struct{ field, want string }{
		{"design_id", "d1"},
		{"section_id", "s2"},
		{"name", "main"},
	} {
		got := msg.Get(md.Fields().ByName(protoreflect.Name(tc.field))).String()
		if got != tc.want {
			t.Errorf("%s = %q, want %q", tc.field, got, tc.want)
		}
	}
}

// TestPopulateEmptyMap verifies that an empty params map is a no-op.
func TestPopulateEmptyMap(t *testing.T) {
	_, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "id", Number: 1, TypeName: "string"},
		},
	}})

	err := fields.PopulateFromMap(msg, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestPopulateUnknownField verifies that referencing a non-existent field
// returns an error.
func TestPopulateUnknownField(t *testing.T) {
	_, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "id", Number: 1, TypeName: "string"},
		},
	}})

	err := fields.PopulateFieldFromPath(msg, "nonexistent", "val")
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

// TestPopulateInvalidCoercion verifies that an unparseable value returns
// an error (e.g., "abc" for an int32 field).
func TestPopulateInvalidCoercion(t *testing.T) {
	_, msg := buildDynMsg(t, []testutil.TestMessage{{
		Name: "Req",
		Fields: []testutil.TestField{
			{Name: "count", Number: 1, TypeName: "int32"},
		},
	}})

	err := fields.PopulateFieldFromPath(msg, "count", "not-a-number")
	if err == nil {
		t.Fatal("expected error for invalid int32 value")
	}
}

// TestPopulateEnumByName verifies setting an enum field by its name.
func TestPopulateEnumByName(t *testing.T) {
	plugin := testutil.CreateTestPlugin(t, &testutil.TestProtoSet{
		Files: []testutil.TestFile{{
			Name: "test.proto",
			Pkg:  "test.v1",
			Enums: []testutil.TestEnum{{
				Name: "Status",
				Values: []testutil.TestEnumValue{
					{Name: "UNKNOWN", Number: 0},
					{Name: "ACTIVE", Number: 1},
					{Name: "INACTIVE", Number: 2},
				},
			}},
			Messages: []testutil.TestMessage{{
				Name: "Req",
				Fields: []testutil.TestField{
					{Name: "status", Number: 1, EnumType: "test.v1.Status"},
				},
			}},
		}},
	})

	var reqDesc protoreflect.MessageDescriptor
	for _, f := range plugin.Files {
		for _, m := range f.Messages {
			if string(m.Desc.Name()) == "Req" {
				reqDesc = m.Desc
			}
		}
	}

	msg := dynamicpb.NewMessage(reqDesc)
	if err := fields.PopulateFieldFromPath(msg, "status", "ACTIVE"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fd := reqDesc.Fields().ByName("status")
	got := msg.Get(fd).Enum()
	if got != 1 {
		t.Errorf("status = %d, want 1 (ACTIVE)", got)
	}
}

// TestPopulateEnumByNumber verifies setting an enum field by numeric value.
func TestPopulateEnumByNumber(t *testing.T) {
	plugin := testutil.CreateTestPlugin(t, &testutil.TestProtoSet{
		Files: []testutil.TestFile{{
			Name: "test.proto",
			Pkg:  "test.v1",
			Enums: []testutil.TestEnum{{
				Name: "Status",
				Values: []testutil.TestEnumValue{
					{Name: "UNKNOWN", Number: 0},
					{Name: "ACTIVE", Number: 1},
				},
			}},
			Messages: []testutil.TestMessage{{
				Name: "Req",
				Fields: []testutil.TestField{
					{Name: "status", Number: 1, EnumType: "test.v1.Status"},
				},
			}},
		}},
	})

	var reqDesc protoreflect.MessageDescriptor
	for _, f := range plugin.Files {
		for _, m := range f.Messages {
			if string(m.Desc.Name()) == "Req" {
				reqDesc = m.Desc
			}
		}
	}

	msg := dynamicpb.NewMessage(reqDesc)
	if err := fields.PopulateFieldFromPath(msg, "status", "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fd := reqDesc.Fields().ByName("status")
	got := msg.Get(fd).Enum()
	if got != 1 {
		t.Errorf("status = %d, want 1 (ACTIVE)", got)
	}
}

