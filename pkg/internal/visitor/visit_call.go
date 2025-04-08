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
	"go/ast"
	"go/types"
	"strings"

	"fillmore-labs.com/zerolint/pkg/internal/checker"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
	"golang.org/x/tools/go/analysis"
)

// visitCall checks for type casts T(x), errors.Is(x, y), errors.As(x, y) and new(T).
func (v *Visitor) visitCall(n *ast.CallExpr) bool { //nolint:cyclop
	switch funType := v.check.TypesInfo().Types[n.Fun]; {
	case funType.IsBuiltin(): // Check for calls to new(T).
		if v.level.Below(level.Extended) {
			return true
		}

		return v.visitBuiltin(n)

	case funType.IsType(): // Check for type casts T(...).
		if v.level.Below(level.Extended) {
			return true
		}

		return v.visitCast(n, funType.Type)

	case funType.IsValue():
		if id, ok := n.Fun.(*ast.Ident); ok && isCfunc(id.Name) {
			return false
		}

		if v.level.AtLeast(level.Full) {
			if sig, ok := funType.Type.(*types.Signature); ok {
				v.visitCallArgs(sig, n.Args)
			}
		}

		// Check for errors.Is(x, y) and errors.As(x, y).
		return v.visitCallFun(n)

	default:
		return true
	}
}

func isCfunc(name string) bool {
	return strings.HasPrefix(name, "_Cfunc_")
}

// visitCallArgs checks explicit nil arguments to pointers to zero-sized parameters.
func (v *Visitor) visitCallArgs(sig *types.Signature, args []ast.Expr) {
	params := sig.Params()
	for i := range min(params.Len(), len(args)) { // Don't deal with variadic functions for now
		param := params.At(i)
		if elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(param.Type()); zeroSized {
			arg := args[i]
			if tv, ok := v.check.TypesInfo().Types[arg]; ok && tv.IsNil() {
				var cM checker.CategorizedMessage
				if name := param.Name(); name == "" {
					cM = msgFormatf(catArgumentNil, valueMethod, "passing explicit nil as parameter of type %q", elem)
				} else {
					cM = msgFormatf(catArgumentNil, valueMethod,
						"passing explicit nil as parameter %q pointing to zero-sized type %q", name, elem)
				}

				fixes := v.check.ReplaceWithZeroValue(arg, elem)
				v.check.Report(arg, cM, fixes)
			}
		}
	}
}

// visitCallFun checks for encoding/json#Decoder.Decode, json.Unmarshal, errors.Is and errors.As.
func (v *Visitor) visitCallFun(n *ast.CallExpr) bool {
	switch fun := ast.Unparen(n.Fun).(type) {
	case *ast.SelectorExpr:
		if sel, ok := v.check.TypesInfo().Selections[fun]; ok {
			// Selection expression
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
// For method values, it delegates to visitMethodVal.
// For method expressions with receivers that are pointers to zero-sized types, it reports an issue.
func (v *Visitor) visitCallSelection(fun *ast.SelectorExpr, sel *types.Selection) bool {
	switch sel.Kind() { //nolint:exhaustive
	case types.MethodVal:
		// Delegate selections like encoding/json#Decoder.Decode.
		return visitMethodVal(sel)

	case types.MethodExpr:
		// Method used as a function value (e.g., (*T).Method).
		if v.level.Below(level.Extended) {
			return true
		}

		if elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(sel.Recv()); zeroSized {
			cM := msgFormatf(catMethodExpression, valueMethod,
				"method expression receiver is pointer to zero-size type %q", elem)
			fixes := v.check.RemoveStar(fun.X)
			v.check.Report(fun, cM, fixes)
		}

		return true

	default:
		return true
	}
}

// visitCallIdent processes encoding/json.Unmarshal (ignored as it requires pointer arguments),
// errors.Is and errors.As from the standard library or golang.org/x/exp/errors.
func (v *Visitor) visitCallIdent(n *ast.CallExpr, fun *ast.Ident) bool { //nolint:cyclop
	obj := v.check.TypesInfo().Uses[fun]
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
			return true // We are only interested in comparisons.
		}

		switch name {
		case "As":
			return false // Do not report pointers in errors.As(..., ...).

		case "Is":
			return v.visitCmp(n, n.Args[0], n.Args[1]) // Delegate analysis of errors.Is(..., ...) to visitCmp.

		default:
			return true
		}

	case "github.com/stretchr/testify/assert", "github.com/stretchr/testify/require":
		if len(n.Args) < 3 { //nolint:mnd
			return true // We are only interested in comparisons.
		}

		switch name {
		case "ErrorAs", "ErrorAsf", "NotErrorAs", "NotErrorAsf":
			return false // Do not report pointers in ....ErrorAs(t, ..., ...).

		case "ErrorIs", "ErrorIsf", "NotErrorIs", "NotErrorIsf":
			return v.visitCmp(n, n.Args[1], n.Args[2]) // Delegate analysis of ErrorIs(t, ..., ...) to visitCmp.

		case "Equal", "NotEqual":
			return true // assert.Equal does not compare object identity, but uses [reflect.DeepEqual].

		default:
			return true
		}

	default:
		return true
	}
}

// visitMethodVal checks for method calls on specific receivers, particularly
// looking for the Decode method on json.Decoder, which requires pointers.
func visitMethodVal(sel *types.Selection) bool {
	fun, ok := sel.Obj().(*types.Func)
	if !ok || fun.Name() != "Decode" { // We are only interested in `Decode`.
		return true
	}

	recv := fun.Signature().Recv().Type() // I'm not using sel.Recv(), since [json.Decoder] could be embedded

	typeName, ok := pointerToTypeName(recv)
	if !ok {
		return true
	}

	// Check for method receiver *encoding/json.Decoder.
	if typeName.Pkg().Path() == "encoding/json" && typeName.Name() == "Decoder" {
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

	elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(t)
	if !zeroSized { // Not a pointer to a zero-sized type.
		return true
	}

	tv, ok := v.check.TypesInfo().Types[n.Args[0]]
	if !ok {
		return true
	}

	var (
		cM      checker.CategorizedMessage
		fixes   []analysis.SuggestedFix
		proceed bool
	)

	if tv.IsNil() {
		proceed = false
		cM = msgFormatf(catCastNil, valueMethod, "cast of nil to pointer to zero-size type %q", elem)

		if s, ok := ast.Unparen(n.Fun).(*ast.StarExpr); ok {
			fixes = v.check.MakePure(n, s.X)
		}
	} else {
		proceed = true
		cM = msgFormatf(catCast, valueMethod, "cast of expression of type %q to pointer to zero-size type %q", tv.Type, elem)
		fixes = v.check.RemoveStar(n.Fun)
	}

	v.check.Report(n, cM, fixes)

	return proceed
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
	argType := v.check.TypesInfo().TypeOf(arg)

	valueMethod, zeroSized := v.check.ZeroSizedType(argType)
	if !zeroSized {
		return true
	}

	cM := msgFormatf(catNew, valueMethod, "new called on zero-sized type %q", argType)
	fixes := v.check.MakePure(n, arg)
	v.check.Report(n, cM, fixes)

	return len(fixes) == 0
}
