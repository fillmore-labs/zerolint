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

	"fillmore-labs.com/zerolint/pkg/internal/typeutil"
)

// visitCallFunc processes encoding/json.Unmarshal (ignored as it requires pointer arguments),
// errors.Is and errors.As from the standard library or golang.org/x/exp/errors.
func (v *Visitor) visitCallFunc(n *ast.CallExpr, fun *types.Func, methodExpr bool) bool {
	if len(n.Args) == 0 { // Plain function call
		return true
	}

	funcName := typeutil.NewFuncName(fun)
	if ftyp, ok := functions[funcName]; ok {
		base := 0
		if methodExpr {
			base = 1
		}

		switch ftyp {
		case funcDecode:
			return false // Do not report pointers in json.Unmarshal(..., ...).

		case funcCmp0:
			if len(n.Args) < base+2 { // Multi-valued argument
				return true
			}

			return v.visitCmp(n, n.Args[base], n.Args[base+1]) // Delegate analysis of errors.Is(..., ...) to visitCmp.

		case funcCmp1:
			if len(n.Args) < base+3 { // Multi-valued argument
				return true
			}

			return v.visitCmp(n, n.Args[base+1], n.Args[base+2]) // Delegate analysis of ErrorIs(t, ..., ...) to visitCmp.

		case funcNone: // should not happen
			v.Diag.LogErrorf(n, "Unconfigured function %s", funcName)

			return true
		}
	}

	return true
}
