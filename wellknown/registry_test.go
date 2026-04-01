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

import "testing"

func TestRegisterAndGet(t *testing.T) {
	r := NewRegistry()
	r.Register("google.protobuf.Timestamp", "time.Time", "time", true)

	m, ok := r.Get("google.protobuf.Timestamp")
	if !ok {
		t.Fatal("expected to find Timestamp mapping")
	}
	if m.TargetType != "time.Time" {
		t.Errorf("expected TargetType 'time.Time', got %q", m.TargetType)
	}
	if m.ImportPath != "time" {
		t.Errorf("expected ImportPath 'time', got %q", m.ImportPath)
	}
	if !m.IsNative {
		t.Error("expected IsNative to be true")
	}
}

func TestGetMissing(t *testing.T) {
	r := NewRegistry()
	_, ok := r.Get("google.protobuf.Nonexistent")
	if ok {
		t.Error("expected not to find nonexistent mapping")
	}
}

func TestIsWellKnown(t *testing.T) {
	r := NewRegistry()
	r.Register("google.protobuf.Duration", "Duration", "google/protobuf", false)

	if !r.IsWellKnown("google.protobuf.Duration") {
		t.Error("expected Duration to be well-known")
	}
	if r.IsWellKnown("google.protobuf.Unknown") {
		t.Error("expected Unknown to not be well-known")
	}
}

func TestAllMappings(t *testing.T) {
	r := NewRegistry()
	r.Register("google.protobuf.Timestamp", "time.Time", "time", true)
	r.Register("google.protobuf.Duration", "time.Duration", "time", true)

	all := r.AllMappings()
	if len(all) != 2 {
		t.Fatalf("expected 2 mappings, got %d", len(all))
	}

	// Ensure it's a copy — modifying returned map shouldn't affect registry
	delete(all, "google.protobuf.Timestamp")
	if !r.IsWellKnown("google.protobuf.Timestamp") {
		t.Error("deleting from AllMappings result should not affect registry")
	}
}

func TestRegisterOverwrite(t *testing.T) {
	r := NewRegistry()
	r.Register("google.protobuf.Timestamp", "string", "", false)
	r.Register("google.protobuf.Timestamp", "time.Time", "time", true)

	m, ok := r.Get("google.protobuf.Timestamp")
	if !ok {
		t.Fatal("expected to find Timestamp mapping")
	}
	if m.TargetType != "time.Time" {
		t.Errorf("expected overwritten TargetType 'time.Time', got %q", m.TargetType)
	}
}

func TestWellKnownProtoTypes(t *testing.T) {
	types := WellKnownProtoTypes()
	if len(types) != 24 {
		t.Errorf("expected 24 well-known types, got %d", len(types))
	}
	// Spot check a few
	found := make(map[string]bool)
	for _, typ := range types {
		found[typ] = true
	}
	expected := []string{
		"google.protobuf.Timestamp",
		"google.protobuf.Duration",
		"google.protobuf.Any",
		"google.protobuf.FieldMask",
		"google.protobuf.SourceContext",
	}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected to find %q in WellKnownProtoTypes", e)
		}
	}
}
