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

func TestDiag_Qualifier(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name              string
		currentPkgPath    string
		currentFileSrc    string
		targetPkgProvider func(currentPkg *types.Package) *types.Package
		setupDiag         func(r *Diag, f *ast.File) // Optional setup, e.g., for malformed imports
		needsImport       bool
		want              string
	}{
		{
			name:           "target package is nil",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return nil
			},
			want: "",
		},
		{
			name:           "target package is current package",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main`,
			targetPkgProvider: func(currentPkg *types.Package) *types.Package {
				return currentPkg
			},
			want: "",
		},
		{
			name:           "target not imported (no imports in file)",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			needsImport: true,
			want:        "\"example.com/other\"",
		},
		{
			name:           "target imported without alias",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import "example.com/other"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "otherpkgname") // Name() should be used
			},
			want: "otherpkgname",
		},
		{
			name:           "target imported with alias",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import custom "example.com/other"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			want: "custom",
		},
		{
			name:           "target imported with dot import",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import . "example.com/other"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			want: "",
		},
		{
			name:           "target imported with underscore import",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import _ "example.com/other"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			needsImport: true,
			want:        "\"example.com/other\"",
		},
		{
			name:           "target not imported",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main; import "example.com/another"`,
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			needsImport: true,
			want:        "\"example.com/other\"",
		},
		{
			name:           "import path unquote error skips spec",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main`, // AST will be modified
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("malformed/path", "skipped")
			},
			setupDiag: func(_ *Diag, f *ast.File) {
				// Manually add an import spec with a path that strconv.Unquote would fail on if it wasn't already quoted.
				// Forcing a direct error from strconv.Unquote is tricky as parser usually creates valid string literals.
				// This simulates a scenario where i.Path.Value is not a valid quoted string.
				// A real AST from a parser would have `Value: "\"malformed/path\""`.
				// If `Value` was just `"malformed/path"` (no quotes), Unquote returns it as is with an error.
				f.Imports = append(f.Imports, &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: "malformed/path"}})
			},
			needsImport: true,
			want:        "\"malformed/path\"", // Falls through to strconv.Quote(pkg.Path()) as the malformed import is skipped
		},
		{
			name:           "current file is nil",
			currentPkgPath: "example.com/main",
			currentFileSrc: `package main`, // Source doesn't matter here
			targetPkgProvider: func(_ *types.Package) *types.Package {
				return types.NewPackage("example.com/other", "other")
			},
			setupDiag: func(d *Diag, _ *ast.File) {
				d.CurrentFile = nil
			},
			needsImport: true,
			want:        "\"example.com/other\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			currentPkg, fset, currentASTFile := parseSourceImports(t, tt.currentPkgPath, tt.currentFileSrc)
			info := &types.Info{} // Minimal info, not used by Qualifier method

			d := newTestDiag(t, info, currentPkg, fset, currentASTFile)
			if tt.setupDiag != nil {
				tt.setupDiag(d, currentASTFile)
			}

			targetPkgToQualify := tt.targetPkgProvider(currentPkg)
			needsImport := false
			got := d.Qualifier(&needsImport)(targetPkgToQualify)

			if needsImport != tt.needsImport {
				t.Errorf("Qualifier(%s) needsImport %t, want %t", targetPkgToQualify.Path(), needsImport, tt.needsImport)
			}

			if got != tt.want {
				t.Errorf("Qualifier(%s) = %q, want %q", targetPkgToQualify.Path(), got, tt.want)
			}
		})
	}
}
