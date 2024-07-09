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

func (v Visitor) visitCallBasic(x *ast.CallExpr) bool { // only [errors.Is]
	if len(x.Args) != 2 { //nolint:mnd
		return true
	}
	if s, ok := unwrap(x.Fun).(*ast.SelectorExpr); !ok || !v.isErrorsIs(s) {
		return true
	}

	return v.visitCmp(x, x.Args[0], x.Args[1])
}

func (v Visitor) visitCall(x *ast.CallExpr) bool {
	tF := v.TypesInfo.Types[x.Fun]
	if tF.IsType() { // check for type casts
		return v.visitCast(x, tF.Type)
	}

	switch f := unwrap(x.Fun).(type) {
	case *ast.SelectorExpr: // check for [errors.Is]
		if len(x.Args) != 2 || !v.isErrorsIs(f) {
			return true
		}

		return v.visitCmp(x, x.Args[0], x.Args[1])

	case *ast.Ident: // check for [builtin.new]
		if !tF.IsBuiltin() || f.Name != "new" || len(x.Args) != 1 {
			return true
		}

		a := x.Args[0]
		tA := v.TypesInfo.Types[a].Type
		if v.isZeroSizeType(tA) {
			message := fmt.Sprintf("new called on zero-size type %q", tA.String())
			fixes := makePure(x, a)
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
	path, ok := v.pkgImportPath(f.X)

	return ok && (path == "errors" || path == "golang.org/x/exp/errors")
}

func (v Visitor) pkgImportPath(x ast.Expr) (path string, ok bool) {
	if id, ok1 := x.(*ast.Ident); ok1 {
		if pkg, ok2 := v.TypesInfo.Uses[id].(*types.PkgName); ok2 {
			return pkg.Imported().Path(), true
		}
	}

	return "", false
}

func (v Visitor) visitCast(x *ast.CallExpr, t types.Type) bool {
	if len(x.Args) != 1 || !v.TypesInfo.Types[x.Args[0]].IsNil() {
		return true
	}

	p, ok := t.Underlying().(*types.Pointer)
	if !ok {
		return true
	}

	e := p.Elem()
	if !v.isZeroSizeType(e) {
		return true
	}

	var fixes []analysis.SuggestedFix
	if s, ok2 := unwrap(x.Fun).(*ast.StarExpr); ok2 {
		fixes = makePure(x, s.X)
	}

	message := fmt.Sprintf("cast of nil to pointer to zero-size variable of type %q", e.String())
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
