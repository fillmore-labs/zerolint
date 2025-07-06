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

	"fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// visitStar analyzes star expressions (*x).
func (v *Visitor) visitStar(n *ast.StarExpr) bool {
	if v.starSeen(n) {
		return false
	}

	t := v.Diag.TypesInfo().TypeOf(n.X)
	if t == nil { // should not happen
		v.Diag.LogErrorf(n, "Can't find star type")

		return true
	}

	if elem, valueMethod, zeroSized := v.Check.ZeroSizedTypePointer(t); zeroSized {
		// *... where the type of ... is a pointer to a zero-size variable.
		// p := &struct{}{}; _ = *p
		cM := msg.Formatf(msg.CatDeref, valueMethod, "dereferencing pointer to zero-size variable of type %q", elem)
		fixes := v.Diag.RemoveOp(n, n.X)
		v.Diag.Report(n, cM, fixes)

		return len(fixes) == 0
	}

	if v.Level.Below(level.Full) {
		return true
	}

	if valueMethod, zeroSized := v.Check.ZeroSizedType(t); zeroSized {
		// *... where ... is a zero-sized type.
		// type T struct{}; map[]*T, []*T, f[*T]
		cM := msg.Formatf(msg.CatStarType, valueMethod, "pointer to zero-sized type %q", t)
		fixes := v.Diag.RemoveOp(n, n.X)
		v.Diag.Report(n, cM, fixes)

		return len(fixes) == 0
	}

	return true
}
