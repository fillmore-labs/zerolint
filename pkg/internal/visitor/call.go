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
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// visitCallValue checks for encoding/json#Decoder.Decode, json.Unmarshal, errors.Is and errors.As.
func (v Visitor) visitCallBasic(x *ast.CallExpr) bool { //nolint:cyclop
	fun, ok := ast.Unparen(x.Fun).(*ast.SelectorExpr)
	if !ok {
		return true
	}

	if sel, ok := v.Pass.TypesInfo.Selections[fun]; ok {
		// Delegate selections
		return visitSelectionCall(sel)
	}

	// qualified identifiers

	pkg, ok := pkgName(v.Pass, fun)
	if !ok {
		return true
	}

	switch path := pkg.Imported().Path(); {
	case path == "encoding/json" && fun.Sel.Name == "Unmarshal":
		return false // Do not report pointers in json.Unmarshal(..., ...).

	case len(x.Args) != 2 || path != "errors" && path != "golang.org/x/exp/errors":
		return true

	case fun.Sel.Name == "As":
		return false // Do not report pointers in errors.As(..., ...)

	case fun.Sel.Name == "Is":
		// Delegate errors.Is(..., ...) to visitCmp for further analysis of the comparison.
		return v.visitCmp(x, x.Args[0], x.Args[1])

	default:
		return true
	}
}

func visitSelectionCall(sel *types.Selection) bool {
	if sel.Kind() != types.MethodVal {
		return true
	}

	if sel.Recv().String() == "*encoding/json.Decoder" && sel.Obj().Name() == "Decode" {
		return false // Do not report pointers in json.Decoder#Decode
	}

	return true
}

func pkgName(pass *analysis.Pass, fun *ast.SelectorExpr) (*types.PkgName, bool) {
	id, ok := fun.X.(*ast.Ident)
	if !ok {
		return nil, false
	}

	pkg, ok := pass.TypesInfo.Uses[id].(*types.PkgName)
	if !ok {
		return nil, false
	}

	return pkg, true
}

// visitCall checks for type casts T(x), errors.Is(x, y), errors.As(x, y) and new(T).
func (v Visitor) visitCall(x *ast.CallExpr) bool {
	switch funType := v.Pass.TypesInfo.Types[x.Fun]; {
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

// visitCast is called for type casts T(nil).
func (v Visitor) visitCast(t types.Type, x *ast.CallExpr) bool {
	if len(x.Args) != 1 {
		return true
	}
	if n, ok := v.Pass.TypesInfo.Types[x.Args[0]]; !ok || !n.IsNil() {
		return true
	}

	elem, ok := v.zeroSizedTypePointer(t)
	if !ok { // Not a pointer to a zero-sized type.
		return true
	}

	message := fmt.Sprintf("cast of nil to pointer to zero-size variable of type %q", elem)
	var fixes []analysis.SuggestedFix
	if s, ok := ast.Unparen(x.Fun).(*ast.StarExpr); ok {
		fixes = v.makePure(x, s.X)
	}

	v.report(x, message, fixes)

	return fixes == nil
}

// visitCast is called for calls to new(T).
func (v Visitor) visitNew(x *ast.CallExpr) bool {
	if len(x.Args) != 1 {
		return true
	}
	fun, ok := ast.Unparen(x.Fun).(*ast.Ident)
	if !ok || fun.Name != "new" {
		return true
	}

	arg := x.Args[0] // new(arg).
	argType := v.Pass.TypesInfo.TypeOf(arg)
	if !v.zeroSizedType(argType) {
		return true
	}

	message := fmt.Sprintf("new called on zero-sized type %q", argType)
	fixes := v.makePure(x, arg)
	v.report(x, message, fixes)

	return false
}
