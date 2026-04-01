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
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// PackageInfo contains extracted information about a proto message's Go package.
type PackageInfo struct {
	ImportPath string // Clean import path without ;packagename suffix
	Alias      string // Package alias for import statements
}

// ExtractPackageInfo extracts clean import path and alias from a protogen message.
// Handles buf-managed packages with ;packagename suffix.
func ExtractPackageInfo(msg *protogen.Message) PackageInfo {
	if msg == nil {
		return PackageInfo{}
	}
	raw := string(msg.GoIdent.GoImportPath)
	importPath := raw
	if idx := strings.Index(raw, ";"); idx >= 0 {
		importPath = raw[:idx]
	}
	return PackageInfo{
		ImportPath: importPath,
		Alias:      GetPackageAlias(importPath),
	}
}

// ExtractGoPackageName extracts the Go package name from a protogen message.
// Prefers the ;suffix override, falls back to last path segment.
func ExtractGoPackageName(msg *protogen.Message) string {
	if msg == nil {
		return ""
	}
	raw := string(msg.GoIdent.GoImportPath)
	if idx := strings.Index(raw, ";"); idx >= 0 {
		return raw[idx+1:]
	}
	if idx := strings.LastIndex(raw, "/"); idx >= 0 {
		return raw[idx+1:]
	}
	return raw
}

// GetPackageAlias returns the last segment of a package path as an alias.
// "github.com/example/api/v1" -> "v1"
func GetPackageAlias(pkgPath string) string {
	if idx := strings.LastIndex(pkgPath, "/"); idx >= 0 {
		return pkgPath[idx+1:]
	}
	return pkgPath
}
