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

import (
	"path/filepath"
	"strings"
)

// CalculateRelativePath calculates the relative path from one directory to another.
// Uses filepath.Rel with forward-slash normalization and ./ prefix for relative imports.
func CalculateRelativePath(fromPath, toPath string) string {
	fromPath = filepath.Clean(fromPath)
	toPath = filepath.Clean(toPath)
	if fromPath == toPath {
		return "."
	}
	rel, err := filepath.Rel(fromPath, toPath)
	if err != nil {
		return toPath
	}
	rel = filepath.ToSlash(rel)
	if !strings.HasPrefix(rel, ".") && !strings.HasPrefix(rel, "/") {
		rel = "./" + rel
	}
	return rel
}

// BuildPackagePath converts a dot-separated package name to a directory path.
// "library.v1" -> "library/v1"
func BuildPackagePath(packageName string) string {
	return strings.ReplaceAll(packageName, ".", "/")
}

// NormalizePath normalizes a file path (resolves ./.., forward slashes, preserves ./ prefix).
func NormalizePath(path string) string {
	hadDotSlash := strings.HasPrefix(path, "./")
	cleaned := filepath.Clean(path)
	cleaned = filepath.ToSlash(cleaned)
	if hadDotSlash && !strings.HasPrefix(cleaned, ".") && !strings.HasPrefix(cleaned, "/") {
		cleaned = "./" + cleaned
	}
	return cleaned
}

// IsAbsolutePath checks if a path is absolute.
func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// JoinPaths joins path components with normalization, preserving ./ prefix of first component.
func JoinPaths(components ...string) string {
	if len(components) == 0 {
		return ""
	}
	hadDotSlash := strings.HasPrefix(components[0], "./")
	filtered := make([]string, 0, len(components))
	for _, c := range components {
		if c != "" {
			filtered = append(filtered, c)
		}
	}
	if len(filtered) == 0 {
		return ""
	}
	joined := filepath.Join(filtered...)
	result := NormalizePath(joined)
	if hadDotSlash && !strings.HasPrefix(result, ".") && !strings.HasPrefix(result, "/") {
		result = "./" + result
	}
	return result
}
