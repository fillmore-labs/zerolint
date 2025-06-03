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

	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
)

// visitCast checks for type casts:
// - nil to pointers of zero-sized types, like (*struct{})(nil).
// - casts of pointers of zero-sized types to [unsafe.Pointer], like unsafe.Pointer(&struct{}{}).
func (v *visitor) visitCast(n *ast.CallExpr, t types.Type) bool {
	if len(n.Args) != 1 { // should not happen
		v.diag.LogErrorf(n, "Expected one argument, got %d", len(n.Args))

		return true
	}

	arg := n.Args[0]

	// Check for unsafe.Pointer(arg)
	if b, ok := t.(*types.Basic); ok && b.Kind() == types.UnsafePointer {
		tv, ok := v.diag.TypesInfo().Types[arg]
		if !ok { // should not happen
			v.diag.LogErrorf(arg, "Can't find unsafe cast type")

			return true
		}

		elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(tv.Type)
		if !zeroSized {
			return true // Not a pointer to a zero-sized type.
		}

		cM := msg.Formatf(msg.CatCastUnsafe, valueMethod, "cast of pointer to zero-size type %q to unsafe.Pointer", elem)
		v.diag.Report(n, cM, nil)

		return false
	}

	// Check for t(arg) with t being a pointer to a zero-sized type
	elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(t)
	if !zeroSized {
		return true // Not a pointer to a zero-sized type.
	}

	tv, ok := v.diag.TypesInfo().Types[arg]
	if !ok { // should not happen
		v.diag.LogErrorf(arg, "Can't find cast type")

		return true
	}

	if tv.IsNil() {
		cM := msg.Formatf(msg.CatCastNil, valueMethod, "cast of nil to pointer to zero-size type %q", elem)

		var fixes []analysis.SuggestedFix
		if s, ok := ast.Unparen(n.Fun).(*ast.StarExpr); ok {
			fixes = v.diag.MakePure(n, s.X)
		}

		v.diag.Report(n, cM, fixes)

		return false // Don't descend into `nil`
	}

	cM := msg.Formatf(msg.CatCast, valueMethod,
		"cast of expression of type %q to pointer to zero-size type %q", tv.Type, elem)
	fixes := v.removeStar(n.Fun)
	v.diag.Report(n, cM, fixes)

	return true // Descend into the argument expression
}
