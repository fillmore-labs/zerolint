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

package checker_test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/checker"
)

// newTestChecker creates a new Checker initialized for testing.
func newTestChecker(tb testing.TB) *Checker {
	tb.Helper()

	c := &Checker{}
	c.Prepare()

	return c
}

func parseFile(tb testing.TB, filename, src string) (*ast.File, *token.FileSet) {
	tb.Helper()

	fset := token.NewFileSet()
	fset.AddFile(filename, -1, len(src))

	f, err := parser.ParseFile(fset, filename, src, parser.SkipObjectResolution)
	if err != nil {
		tb.Fatalf("failed to parse source: %v", err)
	}

	return f, fset
}

// parseSource is a helper to parse source and get [types.Package].
func parseSource(tb testing.TB, filename, src string) *types.Package {
	tb.Helper()

	f, fset := parseFile(tb, filename, src)

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	pkg, err := conf.Check("testpkg", fset, []*ast.File{f}, info)
	if err != nil {
		tb.Fatalf("failed to type check source: %v", err)
	}

	return pkg
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
