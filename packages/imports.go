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

import "sort"

// ImportSpec represents a Go import with optional alias.
type ImportSpec struct {
	Alias string // Optional import alias
	Path  string // Full import path
}

// ImportMap is a deduplicating map of imports keyed by path.
type ImportMap map[string]ImportSpec

// NewImportMap creates an empty ImportMap.
func NewImportMap() ImportMap {
	return make(ImportMap)
}

// Add adds an import, skipping if path already exists.
func (m ImportMap) Add(spec ImportSpec) {
	if _, exists := m[spec.Path]; !exists {
		m[spec.Path] = spec
	}
}

// ToSlice returns imports sorted by path for deterministic output.
func (m ImportMap) ToSlice() []ImportSpec {
	specs := make([]ImportSpec, 0, len(m))
	for _, spec := range m {
		specs = append(specs, spec)
	}
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].Path < specs[j].Path
	})
	return specs
}
