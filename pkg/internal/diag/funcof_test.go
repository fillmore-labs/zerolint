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
	"testing"
)

//nolint:forcetypeassert
func TestDiag_FuncOf(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name         string
		src          string
		findExpr     func(f *ast.File) ast.Expr
		wantFuncName string
		wantOk       bool
	}{
		{
			name: "simple function identifier",
			src: `package testpkg
func myFunc() int { return 0 }
var _ = myFunc()`,
			findExpr: func(f *ast.File) ast.Expr {
				valSpec := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				callExpr := valSpec.Values[0].(*ast.CallExpr)

				return callExpr.Fun // myFunc
			},
			wantFuncName: "myFunc",
			wantOk:       true,
		},
		{
			name: "selector expression (method call)",
			src: `package testpkg
type S struct{}
func (s *S) myMethod() int { return 0 }
var v S
var _ = v.myMethod()`,
			findExpr: func(f *ast.File) ast.Expr {
				valSpec := f.Decls[3].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				callExpr := valSpec.Values[0].(*ast.CallExpr)

				return callExpr.Fun // v.myMethod
			},
			wantFuncName: "myMethod",
			wantOk:       true,
		},
		{
			name: "generic function with one type parameter (IndexExpr)",
			src: `package testpkg
func myGeneric[T any]() int { return 0 }
var _ = myGeneric[int]()`,
			findExpr: func(f *ast.File) ast.Expr {
				valSpec := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				callExpr := valSpec.Values[0].(*ast.CallExpr)

				return callExpr.Fun // myGeneric[int]
			},
			wantFuncName: "myGeneric",
			wantOk:       true,
		},
		{
			name: "generic function with multiple type parameters (IndexListExpr)",
			src: `package testpkg
func myGeneric[T, K any]() int { return 0 }
var _ = myGeneric[int, string]()`,
			findExpr: func(f *ast.File) ast.Expr {
				valSpec := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				callExpr := valSpec.Values[0].(*ast.CallExpr)

				return callExpr.Fun // myGeneric[int, string]
			},
			wantFuncName: "myGeneric",
			wantOk:       true,
		},
		{
			name: "parenthesized function identifier",
			src: `package testpkg
func myFunc() int { return 0 }
var _ = (myFunc)()`,
			findExpr: func(f *ast.File) ast.Expr {
				valSpec := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				callExpr := valSpec.Values[0].(*ast.CallExpr)

				return callExpr.Fun // (myFunc)
			},
			wantFuncName: "myFunc",
			wantOk:       true,
		},
		{
			name: "not a function",
			src: `package testpkg
var myVar int
var _ = myVar`,
			findExpr: func(f *ast.File) ast.Expr {
				valSpec := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)

				return valSpec.Values[0] // myVar
			},
			wantOk: false,
		},
		{
			name: "not an identifier",
			src: `package testpkg
var _ = 1`,
			findExpr: func(f *ast.File) ast.Expr {
				valSpec := f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)

				return valSpec.Values[0] // 1
			},
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info, pkg, fset, astFile := parseSource(t, "test.go", tt.src)
			d := newTestDiag(t, info, pkg, fset, astFile)

			expr := tt.findExpr(astFile)
			fun, ok := d.FuncOf(expr)

			if ok != tt.wantOk {
				t.Fatalf("FuncOf() ok = %v, want %v", ok, tt.wantOk)
			}

			if !tt.wantOk {
				return // Test passed, nothing more to check
			}

			if fun == nil {
				t.Fatal("FuncOf() fun is nil, but ok was true")
			}

			if fun.Name() != tt.wantFuncName {
				t.Errorf("FuncOf() fun.Name() = %q, want %q", fun.Name(), tt.wantFuncName)
			}
		})
	}
}
