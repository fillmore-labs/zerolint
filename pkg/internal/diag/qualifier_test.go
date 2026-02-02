// Copyright 2024-2025 Oliver Eikemeier. All Rights Reserved.
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
//
// SPDX-License-Identifier: Apache-2.0

package diag_test

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/diag"
)

func TestQualifier_Qualifier(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := [...]struct {
		name              string
		currentPkgPath    string
		currentFileSrc    string
		targetPkgProvider func(currentPkg *types.Package) *types.Package
		importsProvider   func(f *ast.File) []*ast.ImportSpec
		needsImport       bool
		want              string
	}{
		{
			name:              "target package is nil",
			currentPkgPath:    "example.com/main",
			currentFileSrc:    `package main`,
			targetPkgProvider: func(_ *types.Package) *types.Package { return nil },
			importsProvider:   func(f *ast.File) []*ast.ImportSpec { return f.Imports },
			want:              "",
		},
		{
			name:              "target package is current package",
			currentPkgPath:    "example.com/main",
			currentFileSrc:    `package main`,
			targetPkgProvider: func(currentPkg *types.Package) *types.Package { return currentPkg },
			importsProvider:   func(f *ast.File) []*ast.ImportSpec { return f.Imports },
			want:              "",
		},
		{
			name:           "target not imported (no imports in file)",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			importsProvider: func(f *ast.File) []*ast.ImportSpec { return f.Imports },
			needsImport:     true,
			want:            "\"example.com/other\"",
		},
		{
			name:           "target imported without alias",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import "example.com/other"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "otherpkgname") // Name() should be used
			},
			importsProvider: func(f *ast.File) []*ast.ImportSpec { return f.Imports },
			want:            "otherpkgname",
		},
		{
			name:           "target imported with alias",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import custom "example.com/other"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			importsProvider: func(f *ast.File) []*ast.ImportSpec { return f.Imports },
			want:            "custom",
		},
		{
			name:           "target imported with dot import",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import . "example.com/other"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			importsProvider: func(f *ast.File) []*ast.ImportSpec { return f.Imports },
			want:            "",
		},
		{
			name:           "target imported with underscore import",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import _ "example.com/other"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			importsProvider: func(f *ast.File) []*ast.ImportSpec { return f.Imports },
			needsImport:     true,
			want:            "\"example.com/other\"",
		},
		{
			name:           "target not imported",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import "example.com/another"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			importsProvider: func(f *ast.File) []*ast.ImportSpec { return f.Imports },
			needsImport:     true,
			want:            "\"example.com/other\"",
		},
		{
			name:           "import path unquote error skips spec",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main`, // AST will be modified
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("malformed/path", "skipped")
			},
			importsProvider: func(f *ast.File) []*ast.ImportSpec {
				// Manually add an import spec with a path that strconv.Unquote would fail on if it wasn't already quoted.
				// Forcing a direct error from strconv.Unquote is tricky as parser usually creates valid string literals.
				// This simulates a scenario where i.Path.Value is not a valid quoted string.
				// A real AST from a parser would have `Value: "\"malformed/path\""`.
				// If `Value` was just `"malformed/path"` (no quotes), Unquote returns it as is with an error.
				f.Imports = append(f.Imports, &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: "malformed/path"}})

				return f.Imports
			},
			needsImport: true,
			want:        "\"malformed/path\"", // Falls through to strconv.Quote(pkg.Path()) as the malformed import is skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			currentPkg, _, currentASTFile := parseSourceImports(t, tt.currentPkgPath, tt.currentFileSrc)

			q := Qualifier{
				Pkg:     currentPkg,
				Imports: tt.importsProvider(currentASTFile),
			}

			targetPkgToQualify := tt.targetPkgProvider(currentPkg)

			var targetPkgPath string
			if targetPkgToQualify == nil {
				targetPkgPath = "<nil>"
			} else {
				targetPkgPath = targetPkgToQualify.Path()
			}

			if got, want := q.Qualifier(targetPkgToQualify), tt.want; got != want {
				t.Errorf("Qualifier(%s) = %q, want %q", targetPkgPath, got, want)
			}

			// Check side effect
			if got, want := q.NeedsImport, tt.needsImport; got != want {
				t.Errorf("Qualifier(%s) needsImport %t, want %t", targetPkgPath, got, want)
			}
		})
	}
}
