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
func (v *Visitor) visitCall(n *ast.CallExpr) bool {
	switch funType := v.Pass.TypesInfo.Types[n.Fun]; {
	case funType.IsBuiltin(): // Check for calls to new(T).
		if !v.full {
			return true
		}

		return v.visitBuiltin(n)

	case funType.IsType(): // Check for type casts T(...).
		if !v.full {
			return true
		}

		return v.visitCast(n, funType.Type)

	case funType.IsValue(): // Check for errors.Is(x, y) and errors.As(x, y).
		return v.visitCallFun(n)

	default:
		return true
	}
}

// visitCallFun checks for encoding/json#Decoder.Decode, json.Unmarshal, errors.Is and errors.As.
func (v *Visitor) visitCallFun(n *ast.CallExpr) bool {
	switch fun := ast.Unparen(n.Fun).(type) {
	case *ast.SelectorExpr:
		if sel, ok := v.Pass.TypesInfo.Selections[fun]; ok {
			// Selection expression
			if !v.full {
				return true
			}

			return v.visitCallSelection(fun, sel)
		}

		return v.visitCallIdent(n, fun.Sel)

	case *ast.Ident:
		return v.visitCallIdent(n, fun)

	default:
		return true
	}
}

func (v *Visitor) visitCallSelection(fun *ast.SelectorExpr, sel *types.Selection) bool {
	switch sel.Kind() { //nolint:exhaustive
	case types.MethodVal:
		// Delegate selections like encoding/json#Decoder.Decode.
		return visitSelectionCall(sel)

	case types.MethodExpr:
		elem, ok := v.zeroSizedTypePointer(sel.Recv())
		if !ok { // Not a pointer receiver or no pointer to a zero-sized type.
			return true
		}

		message := fmt.Sprintf("method expression receiver is pointer to zero-size variable of type %q", elem)
		v.report(fun, message, nil) // Will be fixed by StarExpr.

		return true

	default:
		return true
	}
}

func (v *Visitor) visitCallIdent(n *ast.CallExpr, fun *ast.Ident) bool {
	obj := v.Pass.TypesInfo.Uses[fun]
	pkg := obj.Pkg()
	if pkg == nil {
		return true
	}
	path, name := pkg.Path(), obj.Name()

	switch path {
	case "encoding/json":
		if name == "Unmarshal" {
			return false // Do not report pointers in json.Unmarshal(..., ...).
		}

		return true

	case "errors", "golang.org/x/exp/errors":
		if len(n.Args) != 2 { //nolint:mnd
			return true
		}

		switch name {
		case "As":
			return false // Do not report pointers in errors.As(..., ...).

		case "Is":
			// Delegate errors.Is(..., ...) to visitCmp for further analysis of the comparison.
			return v.visitCmp(n, n.Args[0], n.Args[1])

		default:
			return true
		}

	default:
		return true
	}
}

func visitSelectionCall(sel *types.Selection) bool {
	fun, ok := sel.Obj().(*types.Func)
	if !ok {
		return true
	}

	typeName, ok := pointerToTypeName(fun.Signature().Recv())
	if !ok {
		return true
	}

	// Check for method call "Decode" on receiver *encoding/json.Decoder.
	if fun.Name() == "Decode" && typeName.Pkg().Path() == "encoding/json" && typeName.Name() == "Decoder" {
		return false // Do not report pointers in json.Decoder#Decode.
	}

	return true
}

func pointerToTypeName(v *types.Var) (*types.TypeName, bool) {
	ptr, ok := types.Unalias(v.Type()).(*types.Pointer)
	if !ok {
		return nil, false
	}

	elem, ok := ptr.Elem().(*types.Named)
	if !ok {
		return nil, false
	}

	return elem.Obj(), true
}

// visitCast is called for type casts T(nil).
func (v *Visitor) visitCast(n *ast.CallExpr, t types.Type) bool {
	if len(n.Args) != 1 {
		return true
	}

	if arg, ok := v.Pass.TypesInfo.Types[n.Args[0]]; !ok || !arg.IsNil() {
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

// visitBuiltin is called for calls to new(T).
func (v *Visitor) visitBuiltin(n *ast.CallExpr) bool {
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

	return fixes == nil
}
