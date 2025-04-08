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

	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// visitStar analyzes star expressions (*x).
func (v *Visitor) visitStar(n *ast.StarExpr) bool {
	if v.check.StarSeen(n) {
		return false
	}

	t := v.check.TypesInfo().TypeOf(n.X)
	if t == nil {
		return true
	}

	if elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(t); zeroSized {
		// *... where the type of ... is a pointer to a zero-size variable.
		// p := &struct{}{}; _ = *p
		cM := msgFormatf(catDeref, valueMethod, "dereferencing pointer to zero-size variable of type %q", elem)
		fixes := v.check.RemoveOp(n, n.X)
		v.check.Report(n, cM, fixes)

		return len(fixes) == 0
	}

	if v.level.Below(level.Full) {
		return true
	}

	if valueMethod, zeroSized := v.check.ZeroSizedType(t); zeroSized {
		// *... where ... is a zero-sized type.
		// type T struct{}; map[]*T, []*T, f[*T]
		cM := msgFormatf(catStarType, valueMethod, "pointer to zero-sized type %q", t)
		fixes := v.check.RemoveOp(n, n.X)
		v.check.Report(n, cM, fixes)

		return len(fixes) == 0
	}

	return true
}
