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
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path"
	"testing"

	"golang.org/x/tools/go/analysis"

	. "fillmore-labs.com/zerolint/pkg/internal/diag"
)

// newTestDiag creates a new [Diag] with a Pass initialized for testing.
func newTestDiag(tb testing.TB,
	info *types.Info, pkg *types.Package, fset *token.FileSet, currentFile *ast.File,
) *Diag {
	tb.Helper()

	pass := &analysis.Pass{
		Pkg:       pkg,
		TypesInfo: info,
		Fset:      fset,
		Files:     []*ast.File{currentFile},
		Report:    func(analysis.Diagnostic) {},
	}

	c := New(pass)
	c.CurrentFile = currentFile // Important for [Diag.Qualifier]

	return c
}

// parseSource is a helper to parse source and get [types.Info] and [types.Package].
func parseSource(tb testing.TB, filename, src string) (*types.Info, *types.Package, *token.FileSet, *ast.File) {
	tb.Helper()

	fset := token.NewFileSet()
	fset.AddFile(filename, -1, len(src))

	f, err := parser.ParseFile(fset, filename, src, parser.SkipObjectResolution)
	if err != nil {
		tb.Fatalf("failed to parse source: %v", err)
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	pkg, err := conf.Check("testpkg", fset, []*ast.File{f}, info)
	if err != nil {
		tb.Fatalf("failed to type Check source: %v", err)
	}

	return info, pkg, fset, f
}

// parseSourceImports is a helper to parse source and get a minimal
// types.Package and ast.File, primarily for testing import resolution.
func parseSourceImports(tb testing.TB, pkgPath, src string) (*types.Package, *token.FileSet, *ast.File) {
	tb.Helper()

	fset := token.NewFileSet()

	// We only need imports, not the full AST or type checking for this specific test.
	f, err := parser.ParseFile(fset, "test.go", src, parser.ImportsOnly)
	if err != nil {
		tb.Fatalf("failed to parse source for imports: %v", err)
	}

	var pkgName string
	if f.Name != nil {
		pkgName = f.Name.Name
	} else {
		pkgName = path.Base(pkgPath)
		if pkgName == "." || pkgName == "/" {
			pkgName = "main" // Default if no package decl
		}
	}

	currentPkg := types.NewPackage(pkgPath, pkgName)

	return currentPkg, fset, f
}

// getType is a helper to get a type by name from parsed source.
func getType(tb testing.TB, pkg *types.Package, name string) types.Type {
	tb.Helper()

	obj := pkg.Scope().Lookup(name)
	if obj == nil {
		tb.Fatalf("type %q not found in package %s", name, pkg.Path())
	}

	tn, ok := obj.(*types.TypeName)
	if !ok {
		tb.Fatalf("%q is not a type name: %T", name, obj)
	}

	return tn.Type()
}
