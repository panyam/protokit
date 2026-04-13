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

import (
	"strings"
	"unicode"
)

// ToCamelCase converts PascalCase to camelCase. "FindBooks" -> "findBooks"
func ToCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// ToPascalCase converts camelCase to PascalCase. "findBooks" -> "FindBooks"
func ToPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// ToSnakeCase converts camelCase/PascalCase to snake_case.
// Handles acronyms correctly: "HTMLParser" -> "html_parser", "GetUserByID" -> "get_user_by_id".
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rune(s[i-1])
				// Insert underscore before uppercase if previous char is lowercase,
				// or if next char is lowercase (handles acronym boundaries like "HTMLParser").
				if unicode.IsLower(prev) {
					result.WriteByte('_')
				} else if unicode.IsUpper(prev) && i+1 < len(s) && unicode.IsLower(rune(s[i+1])) {
					result.WriteByte('_')
				}
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// SanitizeIdentifier ensures a string is a valid identifier.
// "user-name" -> "user_name", "123invalid" -> "_23invalid"
func SanitizeIdentifier(name string) string {
	if len(name) == 0 {
		return "identifier"
	}
	var result strings.Builder
	for i, r := range name {
		if i == 0 {
			if unicode.IsLetter(r) || r == '_' {
				result.WriteRune(r)
			} else {
				result.WriteRune('_')
				if unicode.IsDigit(r) {
					// Skip the first digit — it gets replaced by _
				} else {
					// Non-letter, non-digit, non-underscore: just replace with _
				}
			}
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				result.WriteRune(r)
			} else {
				result.WriteRune('_')
			}
		}
	}
	return result.String()
}
