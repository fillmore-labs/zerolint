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

//nolint:forcetypeassert
package diag_test

import (
	"go/ast"
	"go/types"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestDiag_ReplaceWithZeroValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		src             string
		findNodeAndType func(f *ast.File, pkg *types.Package, info *types.Info) (ast.Node, types.Type)
		expectedNewText string
		expectNilFix    bool
	}{
		{
			name: "named struct",
			src:  "package testpkg\ntype MyStruct struct{}\nvar p *MyStruct\nvar _ = p",
			findNodeAndType: func(f *ast.File, pkg *types.Package, _ *types.Info) (ast.Node, types.Type) {
				valSpec := f.Decls[2].(*ast.GenDecl).Specs[0].(*ast.ValueSpec) // var _ = p
				nodeToReplace := valSpec.Values[0]                             // p
				zeroValueType := getType(t, pkg, "MyStruct")

				return nodeToReplace, zeroValueType
			},
			expectedNewText: "MyStruct{}",
		},
		{
			name: "array type",
			src:  "package testpkg\ntype MyArray [0]int\nvar p *MyArray\nvar _ = p",
			findNodeAndType: func(f *ast.File, pkg *types.Package, _ *types.Info) (ast.Node, types.Type) {
				valSpec := f.Decls[2].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				nodeToReplace := valSpec.Values[0]
				zeroValueType := getType(t, pkg, "MyArray")

				return nodeToReplace, zeroValueType
			},
			expectedNewText: "MyArray{}",
		},
		{
			name: "unsupported type (pointer) for zero value literal",
			src:  "package testpkg\ntype MyStruct struct{}\nvar pp **MyStruct\nvar _ = pp",
			findNodeAndType: func(f *ast.File, _ *types.Package, info *types.Info) (ast.Node, types.Type) {
				valSpec := f.Decls[2].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				nodeToReplace := valSpec.Values[0]
				// Attempting to replace 'pp' with zero value of type '*MyStruct'
				zeroValueType := info.TypeOf(nodeToReplace) // type of pp is *MyStruct

				return nodeToReplace, zeroValueType
			},
			expectNilFix: true, // writeType returns false for *types.Pointer
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info, pkg, fset, astFile := parseSource(t, "test.go", tt.src)
			d := newTestDiag(t, info, pkg, fset, astFile)
			nodeToReplace, typeOfZeroValue := tt.findNodeAndType(astFile, pkg, info)

			fixes := d.ReplaceWithZeroValue(nodeToReplace, typeOfZeroValue)
			assertFix(t, fixes, tt.expectNilFix, "replace by zero value", tt.expectedNewText)
		})
	}
}

func TestDiag_RemoveOp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		src             string
		findNodes       func(f *ast.File, info *types.Info) (ast.Node, ast.Expr) // n, x
		expectedNewText string
	}{
		{
			name: "remove star",
			src:  "package testpkg\nvar i int\nvar pi = &i\nvar _ = *pi",
			findNodes: func(f *ast.File, _ *types.Info) (ast.Node, ast.Expr) {
				valSpec := f.Decls[2].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				starExprN := valSpec.Values[0].(*ast.StarExpr) // *pi
				identX := starExprN.X.(*ast.Ident)             // pi

				return starExprN, identX
			},
			expectedNewText: "pi",
		},
		{
			name: "remove ampersand",
			src:  "package testpkg\nvar i int\nvar _ = &i",
			findNodes: func(f *ast.File, _ *types.Info) (ast.Node, ast.Expr) {
				valSpec := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				unaryExprN := valSpec.Values[0].(*ast.UnaryExpr) // &i
				identX := unaryExprN.X.(*ast.Ident)              // i

				return unaryExprN, identX
			},
			expectedNewText: "i",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info, pkg, fset, astFile := parseSource(t, "test.go", tt.src)
			d := newTestDiag(t, info, pkg, fset, astFile)
			nodeN, exprX := tt.findNodes(astFile, info)

			fixes := d.RemoveOp(nodeN, exprX)
			assertFix(t, fixes, false, "remove operator", tt.expectedNewText)
		})
	}
}

func TestDiag_MakePure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		src             string
		findNodes       func(f *ast.File, info *types.Info) (ast.Node, ast.Expr) // n, x (type expr)
		expectedNewText string
	}{
		{
			name: "new(T)",
			src:  "package testpkg\ntype S struct{}\nvar _ = new(S)",
			findNodes: func(f *ast.File, _ *types.Info) (ast.Node, ast.Expr) {
				valSpec := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				callExprN := valSpec.Values[0].(*ast.CallExpr) // new(S)
				typeExprX := callExprN.Args[0].(*ast.Ident)    // S

				return callExprN, typeExprX
			},
			expectedNewText: "S{}",
		},
		{
			name: "(*T)(nil)",
			src:  "package testpkg\ntype S struct{}\nvar _ = (*S)(nil)",
			findNodes: func(f *ast.File, _ *types.Info) (ast.Node, ast.Expr) {
				valSpec := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
				callExprN := valSpec.Values[0].(*ast.CallExpr) // (*S)(nil)
				parenExpr := callExprN.Fun.(*ast.ParenExpr)    // (*S)
				starExpr := parenExpr.X.(*ast.StarExpr)        // *S
				typeExprX := starExpr.X.(*ast.Ident)           // S

				return callExprN, typeExprX
			},
			expectedNewText: "S{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info, pkg, fset, astFile := parseSource(t, "test.go", tt.src)
			d := newTestDiag(t, info, pkg, fset, astFile)

			nodeN, exprX := tt.findNodes(astFile, info)

			fixes := d.MakePure(nodeN, exprX)
			assertFix(t, fixes, false, "change to pure type", tt.expectedNewText)
		})
	}
}

func assertFix(t *testing.T, fixes []analysis.SuggestedFix, expectNil bool, expectedMsg, expectedNewText string) {
	t.Helper()

	if expectNil {
		if fixes != nil {
			t.Errorf("expected nil fixes, got %+v", fixes)
		}

		return
	}

	if fixes == nil {
		t.Fatal("expected non-nil fixes, got nil")
	}

	if len(fixes) != 1 {
		t.Fatalf("expected 1 fix, got %d", len(fixes))
	}

	fix := fixes[0]
	if !strings.Contains(fix.Message, expectedMsg) { // Message might have (zl:...) suffix
		t.Errorf("fix.Message = %q, want to contain %q", fix.Message, expectedMsg)
	}

	if len(fix.TextEdits) != 1 {
		t.Fatalf("expected 1 text edit, got %d", len(fix.TextEdits))
	}

	edit := fix.TextEdits[0]
	if string(edit.NewText) != expectedNewText {
		t.Errorf("edit.NewText = %q, want %q", string(edit.NewText), expectedNewText)
	}
}
