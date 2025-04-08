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

// visitCallSelection handles method expressions selected from a value (dec.Decode(...)) or type
// ((*json.Decoder).Decode).
// For method values, it delegates to visitSelectionCall.
// For method expressions with receivers that are pointers to zero-sized types, it reports an issue.
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

		message := fmt.Sprintf("method expression receiver is pointer to zero-size variable of type %q (ZL13)", elem)
		fixes := v.removeStar(fun.X)
		v.report(fun, message, fixes)

		return true

	default:
		return true
	}
}

// visitCallIdent processes encoding/json.Unmarshal (ignored as it requires pointer arguments),
// errors.Is and errors.As from the standard library or golang.org/x/exp/errors.
func (v *Visitor) visitCallIdent(n *ast.CallExpr, fun *ast.Ident) bool { //nolint:cyclop
	obj := v.Pass.TypesInfo.Uses[fun]
	if obj == nil {
		return true
	}
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

// visitSelectionCall checks for method calls on specific receivers, particularly
// looking for the Decode method on json.Decoder, which requires pointers.
func visitSelectionCall(sel *types.Selection) bool {
	fun, ok := sel.Obj().(*types.Func)
	if !ok {
		return true
	}

	typeName, ok := pointerToTypeName(fun.Signature().Recv().Type())
	if !ok {
		return true
	}

	// Check for method call "Decode" on receiver *encoding/json.Decoder.
	if fun.Name() == "Decode" && typeName.Pkg().Path() == "encoding/json" && typeName.Name() == "Decoder" {
		return false // Do not report pointers in json.Decoder#Decode.
	}

	return true
}

// pointerToTypeName extracts the underlying named type from a pointer type.
func pointerToTypeName(t types.Type) (*types.TypeName, bool) {
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return nil, false
	}

	elem, ok := ptr.Elem().(*types.Named)
	if !ok {
		return nil, false
	}

	return elem.Obj(), true
}

// visitCast checks for type casts of nil to pointers of zero-sized types, like (*struct{})(nil).
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

	message := fmt.Sprintf("cast of nil to pointer to zero-size variable of type %q (ZL11)", elem)
	var fixes []analysis.SuggestedFix
	if s, ok := ast.Unparen(n.Fun).(*ast.StarExpr); ok {
		fixes = v.makePure(n, s.X)
	}

	v.report(n, message, fixes)

	return fixes == nil
}

// visitBuiltin examines calls to new(T), where T is a zero-sized type.
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

	message := fmt.Sprintf("new called on zero-sized type %q (ZL10)", argType)
	fixes := v.makePure(n, arg)
	v.report(n, message, fixes)

	return fixes == nil
}
