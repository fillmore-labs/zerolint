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

func (v Visitor) visitCallBasic(x *ast.CallExpr) bool {
	if len(x.Args) != 2 { //nolint:mnd
		return true
	}
	if s, ok := unwrap(x.Fun).(*ast.SelectorExpr); !ok || !v.isErrorsIs(s) {
		return true
	}

	// errors.Is(..., ...)
	return v.visitCmp(x, x.Args[0], x.Args[1])
}

func (v Visitor) visitCall(x *ast.CallExpr) bool {
	funType := v.TypesInfo.Types[x.Fun]
	if funType.IsType() { // check for type casts
		// (T)(...)
		return v.visitCast(x, funType.Type)
	}

	switch fun := unwrap(x.Fun).(type) {
	case *ast.SelectorExpr:
		if len(x.Args) != 2 || !v.isErrorsIs(fun) {
			return true
		}

		// errors.Is(..., ...)
		return v.visitCmp(x, x.Args[0], x.Args[1])

	case *ast.Ident:
		if !funType.IsBuiltin() || fun.Name != "new" || len(x.Args) != 1 {
			return true
		}

		// new(...)
		arg := x.Args[0]
		argType := v.TypesInfo.Types[arg].Type
		if v.isZeroSizeType(argType) {
			message := fmt.Sprintf("new called on zero-size type %q", argType)
			fixes := makePure(x, arg)
			v.report(x, message, fixes)

			return false
		}
	}

	return true
}

func (v Visitor) isErrorsIs(f *ast.SelectorExpr) bool {
	if f.Sel.Name != "Is" {
		return false
	}

	id, ok := f.X.(*ast.Ident)
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

func (v Visitor) visitCast(x *ast.CallExpr, t types.Type) bool {
	if len(x.Args) != 1 || !v.TypesInfo.Types[x.Args[0]].IsNil() {
		return true
	}

	p, ok := t.Underlying().(*types.Pointer)
	if !ok || !v.isZeroSizeType(p.Elem()) {
		return true
	}

	// (*...)(nil)
	message := fmt.Sprintf("cast of nil to pointer to zero-size variable of type %q", p.Elem())

	var fixes []analysis.SuggestedFix
	if s, ok2 := unwrap(x.Fun).(*ast.StarExpr); ok2 {
		fixes = makePure(x, s.X)
	}

	v.report(x, message, fixes)

	return fixes == nil
}

func unwrap(e ast.Expr) ast.Expr {
	x := e
	for p, ok := x.(*ast.ParenExpr); ok; p, ok = x.(*ast.ParenExpr) {
		x = p.X
	}

	return x
}

func makePure(n ast.Node, x ast.Expr) []analysis.SuggestedFix {
	var fixes []analysis.SuggestedFix
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), x); err == nil {
		buf.WriteString("{}")
		edit := analysis.TextEdit{
			Pos:     n.Pos(),
			End:     n.End(),
			NewText: buf.Bytes(),
		}
		fixes = []analysis.SuggestedFix{
			{Message: "change to pure type", TextEdits: []analysis.TextEdit{edit}},
		}
	}

	return fixes
}
