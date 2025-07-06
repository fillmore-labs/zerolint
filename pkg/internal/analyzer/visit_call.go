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

package analyzer

import (
	"go/ast"
	"go/types"

	"fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
	"fillmore-labs.com/zerolint/pkg/internal/typeutil"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// visitCall analyzes call expressions. It acts as a dispatcher, handling:
//   - Built-in calls, such as new(T).
//   - Type conversions, such as (*T)(nil).
//   - Function and method calls, checking for nil arguments to zero-sized pointer
//     parameters and special-casing functions like errors.Is or json.Unmarshal.
func (v *Visitor) visitCall(n *ast.CallExpr) bool {
	switch funType := v.Diag.TypesInfo().Types[n.Fun]; {
	case funType.IsBuiltin(): // Check for calls to new(T).
		if v.Level.Below(level.Extended) {
			return true
		}

		return v.visitBuiltin(n)

	case funType.IsType(): // Check for type casts T(...).
		if v.Level.Below(level.Extended) {
			return true
		}

		return v.visitCast(n, funType.Type)

	case funType.IsValue():
		// Retrieve the definition of the called function.
		fun, methodExpr, ok := typeutil.FuncOf(v.Diag.TypesInfo(), n.Fun)
		if !ok { // *ast.FuncLit or *types.Var
			return true
		}

		if isCfunc(fun.Name()) {
			return false
		}

		// visitCallFun checks for encoding/json#Decoder.Decode, json.Unmarshal, errors.Is and errors.As.
		if !v.visitCallFunc(n, fun, methodExpr) {
			return false
		}

		if methodExpr && v.Level.AtLeast(level.Extended) {
			if e, ok := ast.Unparen(n.Fun).(*ast.SelectorExpr); ok {
				// Selection expression
				if sel, ok := v.Diag.TypesInfo().Selections[e]; ok && sel.Kind() == types.MethodExpr {
					// Method used as a function value (e.g., (*T).Method).
					v.visitCallMethodExpr(e, sel)
				}
			}
		}

		if v.Level.AtLeast(level.Full) {
			if sig, ok := funType.Type.(*types.Signature); ok {
				v.visitCallArgs(sig, n.Args) // Check for nil arguments to zero-sized pointer parameters.
			}
		}

		return true

	default: // should not happen
		v.Diag.LogErrorf(n, "Expected type or value, got %v", funType)

		return true
	}
}

// visitCallSelection handles method expressions selected from a value (dec.Decode(...)) or type
// ((*json.Decoder).Decode).
// For method expressions with receivers that are pointers to zero-sized types, it reports an issue.
func (v *Visitor) visitCallMethodExpr(n *ast.SelectorExpr, sel *types.Selection) {
	// We care about the receiver of the method expression, *not* the signature of the possibly embedded function:
	//	recv := fun.Signature().Recv().Type()
	recv := sel.Recv()

	if elem, valueMethod, zeroSized := v.Check.ZeroSizedTypePointer(recv); zeroSized {
		cM := msg.Formatf(msg.CatMethodExpression, valueMethod,
			"method expression receiver is pointer to zero-size type %q", elem)
		fixes := v.removeStar(n.X)
		v.Diag.Report(n, cM, fixes)
	}
}
