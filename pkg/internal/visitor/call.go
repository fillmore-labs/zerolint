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

// visitCall checks for type casts T(x), errors.Is(x, y), errors.As(x, y) and new(T).
func (v *visitorInternal) visitCall(n *ast.CallExpr) bool {
	switch funType := v.Pass.TypesInfo.Types[n.Fun]; {
	case funType.IsBuiltin(): // Check for calls to new(T).
		return v.visitNew(n)

	case funType.IsType(): // Check for type casts T(...).
		return v.visitCast(funType.Type, n)

	case funType.IsValue(): // Check for errors.Is(x, y) and errors.As(x, y).
		return v.visitCallBasic(n)

	default:
		return true
	}
}

// visitCallBasic checks for encoding/json#Decoder.Decode, json.Unmarshal, errors.Is and errors.As.
func (v *visitorInternal) visitCallBasic(n *ast.CallExpr) bool { //nolint:cyclop
	sel, ok := ast.Unparen(n.Fun).(*ast.SelectorExpr)
	if !ok {
		return true
	}

	if sel, ok := v.Pass.TypesInfo.Selections[sel]; ok {
		// Delegate selections like encoding/json#Decoder.Decode
		return visitSelectionCall(sel)
	}

	// qualified identifiers
	pkg, ok := pkgName(v.Pass, sel)
	if !ok {
		return true
	}

	switch path := pkg.Imported().Path(); {
	case path == "encoding/json" && sel.Sel.Name == "Unmarshal":
		return false // Do not report pointers in json.Unmarshal(..., ...).

	case len(n.Args) != 2 || path != "errors" && path != "golang.org/x/exp/errors":
		return true

	case sel.Sel.Name == "As":
		return false // Do not report pointers in errors.As(..., ...)

	case sel.Sel.Name == "Is":
		// Delegate errors.Is(..., ...) to visitCmp for further analysis of the comparison.
		return v.visitCmp(n, n.Args[0], n.Args[1])

	default:
		return true
	}
}

func pkgName(pass *analysis.Pass, sel *ast.SelectorExpr) (*types.PkgName, bool) {
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return nil, false
	}

	pkg, ok := pass.TypesInfo.Uses[id].(*types.PkgName)
	if !ok {
		return nil, false
	}

	return pkg, true
}

func visitSelectionCall(sel *types.Selection) bool {
	if sel.Kind() != types.MethodVal {
		return true
	}

	fun, ok := sel.Obj().(*types.Func)
	if !ok {
		return true
	}

	elem, ok := receiverPointerToNamed(fun)
	if !ok {
		return true
	}

	// check for method call "Decode" on receiver *encoding/json.Decoder
	if fun.Name() == "Decode" && elem.Obj().Pkg().Path() == "encoding/json" && elem.Obj().Name() == "Decoder" {
		return false // Do not report pointers in json.Decoder#Decode
	}

	return true
}

func receiverPointerToNamed(fun *types.Func) (*types.Named, bool) {
	recv := fun.Signature().Recv()
	if recv == nil {
		return nil, false
	}

	ptr, ok := types.Unalias(recv.Type()).(*types.Pointer)
	if !ok {
		return nil, false
	}

	elem, ok := ptr.Elem().(*types.Named)
	if !ok {
		return nil, false
	}

	return elem, true
}

// visitCast is called for type casts T(nil).
func (v *visitorInternal) visitCast(t types.Type, n *ast.CallExpr) bool {
	if len(n.Args) != 1 {
		return true
	}
	if n, ok := v.Pass.TypesInfo.Types[n.Args[0]]; !ok || !n.IsNil() {
		return true
	}

	elem, ok := v.zeroSizedTypePointer(t)
	if !ok { // Not a pointer to a zero-sized type.
		return true
	}

	message := fmt.Sprintf("cast of nil to pointer to zero-size variable of type %q", elem)
	var fixes []analysis.SuggestedFix
	if s, ok := ast.Unparen(n.Fun).(*ast.StarExpr); ok {
		fixes = v.makePure(n, s.X)
	}

	v.report(n, message, fixes)

	return fixes == nil
}

// visitNew is called for calls to new(T).
func (v *visitorInternal) visitNew(n *ast.CallExpr) bool {
	if len(n.Args) != 1 {
		return true
	}
	fun, ok := ast.Unparen(n.Fun).(*ast.Ident)
	if !ok || fun.Name != "new" {
		return true
	}

	arg := n.Args[0] // new(arg).
	argType := v.Pass.TypesInfo.TypeOf(arg)
	if !v.zeroSizedType(argType) {
		return true
	}

	message := fmt.Sprintf("new called on zero-sized type %q", argType)
	fixes := v.makePure(n, arg)
	v.report(n, message, fixes)

	return false
}
