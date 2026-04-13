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

package names

import "testing"

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FindBooks", "findBooks"},
		{"GetUser", "getUser"},
		{"ID", "iD"},
		{"a", "a"},
		{"A", "a"},
		{"", ""},
		{"alreadyCamel", "alreadyCamel"},
		{"X", "x"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			if got != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"findBooks", "FindBooks"},
		{"getUser", "GetUser"},
		{"id", "Id"},
		{"a", "A"},
		{"A", "A"},
		{"", ""},
		{"AlreadyPascal", "AlreadyPascal"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToPascalCase(tt.input)
			if got != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FindBooks", "find_books"},
		{"getUser", "get_user"},
		{"ID", "id"},
		{"simpleCase", "simple_case"},
		{"HTTPServer", "http_server"},
		{"HTMLParser", "html_parser"},
		{"GetUserByID", "get_user_by_id"},
		{"OAuth2Token", "o_auth2token"},
		{"a", "a"},
		{"A", "a"},
		{"", ""},
		{"alreadysnake", "alreadysnake"},
		{"ABC", "abc"},
		{"CreateDeckV2", "create_deck_v2"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToSnakeCase(tt.input)
			if got != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSanitizeIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user-name", "user_name"},
		{"123invalid", "_23invalid"},
		{"valid_name", "valid_name"},
		{"_private", "_private"},
		{"has.dot", "has_dot"},
		{"has spaces", "has_spaces"},
		{"", "identifier"},
		{"a", "a"},
		{"CamelCase", "CamelCase"},
		{"with$special@chars", "with_special_chars"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := SanitizeIdentifier(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeIdentifier(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
