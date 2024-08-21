// Copyright 2024 Oliver Eikemeier. All Rights Reserved.
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

package visitor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// visitCallBasic checks for errors.Is(x, y) and errors.As(x, y).
func (v Visitor) visitCallBasic(x *ast.CallExpr) bool {
	if len(x.Args) != 2 { //nolint:mnd
		return true
	}
	fun, ok := unwrap(x.Fun).(*ast.SelectorExpr)
	if !ok || !v.isErrors(fun.X) {
		return true
	}

	if fun.Sel.Name == "As" { // Do not report pointers in errors.As(..., ...).
		return false
	}

	if fun.Sel.Name != "Is" {
		return true
	}

	// Delegate errors.Is(..., ...) to visitCmp for further analysis of the comparison.
	return v.visitCmp(x, x.Args[0], x.Args[1])
}

// visitCall checks for type casts T(x), errors.Is(x, y), errors.As(x, y) and new(T).
func (v Visitor) visitCall(x *ast.CallExpr) bool {
	switch funType := v.TypesInfo.Types[x.Fun]; {
	case funType.IsType(): // Check for type casts T(...).
		return v.visitCast(funType.Type, x)

	case funType.IsBuiltin(): // Check for calls to new(T).
		return v.visitNew(x)

	case funType.IsValue(): // Check for errors.Is(x, y) and errors.As(x, y).
		return v.visitCallBasic(x)

	default:
		return true
	}
}

// isErrors checks whether the expression is a package specifier for errors or golang.org/x/exp/errors.
func (v Visitor) isErrors(x ast.Expr) bool {
	id, ok := x.(*ast.Ident)
	if !ok {
		return false
	}

	pkg, ok := v.TypesInfo.Uses[id].(*types.PkgName)
	if !ok {
		return false
	}

	path := pkg.Imported().Path()

	return path == "errors" || path == "golang.org/x/exp/errors"
}

// visitCast is called for type casts T(nil).
func (v Visitor) visitCast(t types.Type, x *ast.CallExpr) bool {
	if len(x.Args) != 1 || !v.TypesInfo.Types[x.Args[0]].IsNil() {
		return true
	}

	elem, ok := v.zeroSizedTypePointer(t)
	if !ok { // Not a pointer to a zero-sized type.
		return true
	}

	message := fmt.Sprintf("cast of nil to pointer to zero-size variable of type %q", elem)
	var fixes []analysis.SuggestedFix
	if s, ok2 := unwrap(x.Fun).(*ast.StarExpr); ok2 {
		fixes = makePure(x, s.X)
	}

	v.report(x, message, fixes)

	return fixes == nil
}

// visitCast is called for calls to new(T).
func (v Visitor) visitNew(x *ast.CallExpr) bool {
	if len(x.Args) != 1 {
		return true
	}
	fun, ok := unwrap(x.Fun).(*ast.Ident)
	if !ok || fun.Name != "new" {
		return true
	}

	arg := x.Args[0] // new(arg).
	argType := v.TypesInfo.Types[arg].Type
	if !v.isZeroSizedType(argType) {
		return true
	}

	message := fmt.Sprintf("new called on zero-sized type %q", argType)
	fixes := makePure(x, arg)
	v.report(x, message, fixes)

	return false
}

// unwrap removes parentheses from an expression (x).
func unwrap(e ast.Expr) ast.Expr {
	x := e
	for {
		p, ok := x.(*ast.ParenExpr)
		if !ok {
			break
		}
		x = p.X
	}

	return x
}

// makePure adds a suggested fix from (*T)(nil) or new(T) to T{}.
func makePure(n ast.Node, x ast.Expr) []analysis.SuggestedFix {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), x); err != nil {
		return nil
	}
	buf.WriteString("{}")
	edit := analysis.TextEdit{
		Pos:     n.Pos(),
		End:     n.End(),
		NewText: buf.Bytes(),
	}

	return []analysis.SuggestedFix{
		{
			Message:   "change to pure type",
			TextEdits: []analysis.TextEdit{edit},
		},
	}
}
