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

package packages

import "testing"

func TestGetPackageAlias(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"github.com/example/api/v1", "v1"},
		{"github.com/example/pkg", "pkg"},
		{"simple", "simple"},
		{"a/b/c/d", "d"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := GetPackageAlias(tt.input)
			if got != tt.want {
				t.Errorf("GetPackageAlias(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildPackagePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"library.v1", "library/v1"},
		{"com.example.api", "com/example/api"},
		{"single", "single"},
		{"a.b.c.d", "a/b/c/d"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := BuildPackagePath(tt.input)
			if got != tt.want {
				t.Errorf("BuildPackagePath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCalculateRelativePath(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want string
	}{
		{"same directory", "/a/b", "/a/b", "."},
		{"child directory", "/a/b", "/a/b/c", "./c"},
		{"sibling directory", "/a/b", "/a/c", "../c"},
		{"deeper nesting", "/a/b/c", "/a/d/e", "../../d/e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateRelativePath(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("CalculateRelativePath(%q, %q) = %q, want %q", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"./a/b/c", "./a/b/c"},
		{"./a/../b", "./b"},
		{"a/b/c", "a/b/c"},
		{"/a/b/c", "/a/b/c"},
		{"./a/./b", "./a/b"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := NormalizePath(tt.input)
			if got != tt.want {
				t.Errorf("NormalizePath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestImportMap(t *testing.T) {
	m := NewImportMap()
	m.Add(ImportSpec{Path: "fmt"})
	m.Add(ImportSpec{Alias: "v1", Path: "github.com/example/api/v1"})
	m.Add(ImportSpec{Path: "fmt"}) // duplicate, should be skipped

	specs := m.ToSlice()
	if len(specs) != 2 {
		t.Fatalf("expected 2 imports, got %d", len(specs))
	}
	// Should be sorted by path
	if specs[0].Path != "fmt" {
		t.Errorf("expected first import path 'fmt', got %q", specs[0].Path)
	}
	if specs[1].Path != "github.com/example/api/v1" {
		t.Errorf("expected second import path 'github.com/example/api/v1', got %q", specs[1].Path)
	}
}

func TestJoinPaths(t *testing.T) {
	tests := []struct {
		name       string
		components []string
		want       string
	}{
		{"simple join", []string{"a", "b", "c"}, "a/b/c"},
		{"dot-slash preserved", []string{"./a", "b"}, "./a/b"},
		{"empty components filtered", []string{"a", "", "b"}, "a/b"},
		{"no components", []string{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JoinPaths(tt.components...)
			if got != tt.want {
				t.Errorf("JoinPaths(%v) = %q, want %q", tt.components, got, tt.want)
			}
		})
	}
}
