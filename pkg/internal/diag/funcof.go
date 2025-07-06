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

package diag

import (
	"go/ast"
	"go/types"
)

// FuncOf resolves an expression that identifies a function to its underlying [*types.Func] object.
// It handles simple function identifiers, selector expressions (e.g., `pkg.Func`), and
// generic function instantiations (e.g., `Func[T]` or `Func[T, U]`).
// It returns the resolved function and true if successful, or nil and false otherwise.
func (d *Diag) FuncOf(e ast.Expr) (fun *types.Func, ok bool) {
	switch e := ast.Unparen(e).(type) {
	case *ast.Ident:
		if obj, ok := d.pass.TypesInfo.Uses[e]; ok {
			fun, ok = obj.(*types.Func)

			return fun, ok
		}

	case *ast.SelectorExpr:
		// types.Checker calls recordUse from recordSelection
		if obj, ok := d.pass.TypesInfo.Uses[e.Sel]; ok {
			fun, ok = obj.(*types.Func)

			return fun, ok
		}

	case *ast.IndexExpr:
		if typ, ok := d.pass.TypesInfo.Types[e.Index]; ok && typ.IsType() {
			return d.FuncOf(e.X) // Unwrap generic expression
		}

	case *ast.IndexListExpr:
		if len(e.Indices) > 0 {
			if typ, ok := d.pass.TypesInfo.Types[e.Indices[0]]; ok && typ.IsType() {
				return d.FuncOf(e.X) // Unwrap generic expression
			}
		}
	}

	return nil, false
}
