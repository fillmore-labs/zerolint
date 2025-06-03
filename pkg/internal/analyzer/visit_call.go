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

	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// visitCall analyzes call expressions. It acts as a dispatcher, handling:
//   - Built-in calls, such as new(T).
//   - Type conversions, such as (*T)(nil).
//   - Function and method calls, checking for nil arguments to zero-sized pointer
//     parameters and special-casing functions like errors.Is or json.Unmarshal.
func (v *visitor) visitCall(n *ast.CallExpr) bool {
	switch funType := v.diag.TypesInfo().Types[n.Fun]; {
	case funType.IsBuiltin():
		if v.level.Below(level.Extended) {
			return true
		}

		return v.visitBuiltin(n) // Check for calls to new(T).

	case funType.IsType():
		if v.level.Below(level.Extended) {
			return true
		}

		return v.visitCast(n, funType.Type) // Check for type casts T(...).

	case funType.IsValue():
		if !v.visitCallFun(n) { // Check for errors.Is(x, y) and errors.As(x, y).
			return false
		}

		if v.level.AtLeast(level.Full) {
			if sig, ok := funType.Type.(*types.Signature); ok {
				v.visitCallArgs(sig, n.Args) // Check for nil arguments to zero-sized pointer parameters.
			}
		}

		return true

	default: // should not happen
		v.diag.LogErrorf(n, "Expected type or value, got %v", funType)

		return true
	}
}

// visitCallFun checks for encoding/json#Decoder.Decode, json.Unmarshal, errors.Is and errors.As.
func (v *visitor) visitCallFun(n *ast.CallExpr) bool {
	switch fun := ast.Unparen(n.Fun).(type) {
	case *ast.SelectorExpr:
		if sel, ok := v.diag.TypesInfo().Selections[fun]; ok {
			// Selection expression
			if v.level.Below(level.Extended) {
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
