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

package typeutil

import (
	"go/ast"
	"go/types"
)

// FuncOf iteratively unwraps an expression to find the underlying function declaration.
func FuncOf(info *types.Info, ex ast.Expr) (fun *types.Func, methodExpr, ok bool) {
	for {
		switch e := ex.(type) {
		case *ast.Ident:
			fun, ok = info.Uses[e].(*types.Func)

			return fun, false, ok

		case *ast.SelectorExpr:
			sel, ok := info.Selections[e]
			if !ok { // e.Sel is an identifier qualified by e.X
				fun, ok = info.Uses[e.Sel].(*types.Func) // types.Checker calls recordUse for e.Sel from recordSelection.

				return fun, false, ok
			}

			switch sel.Kind() { //nolint:exhaustive
			case types.MethodVal: // e.Sel is a method selector
				fun, ok = sel.Obj().(*types.Func)

				return fun, false, ok

			case types.MethodExpr: // e.Sel is a method expression
				fun, ok = sel.Obj().(*types.Func)

				return fun, true, ok
			}

			return nil, false, false // e.Sel is a struct field selector

		case *ast.IndexExpr: // Generic function instantiation with a type parameter ("myFunc[T]").
			if !info.Types[e.Index].IsType() {
				return nil, false, false // Must be a type parameter, not an array/slice index.
			}

			ex = e.X // Unwrap to the function identifier.

		case *ast.IndexListExpr: // Generic function instantiation with multiple type parameters ("myFunc[T, U]").
			ex = e.X // Unwrap to the function identifier.

		case *ast.ParenExpr: // Parenthesized expression ("(myFunc)")
			ex = e.X // Unwrap to the inner expression.

		default: // The expression does not resolve to a function identifier (could be a function pointer).
			return nil, false, false
		}
	}
}
