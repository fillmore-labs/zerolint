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
)

// visitStar analyzes star expressions (*x).
func (v *Visitor) visitStar(n *ast.StarExpr) bool {
	if v.base.StarSeen(n) {
		return false
	}

	t := v.base.TypesInfo().TypeOf(n.X)
	if cat, message, valueMethod, zeroSized := v.zeroSizedStarExpr(t); zeroSized {
		fixes := v.base.RemoveOp(n, n.X)
		v.base.Report(n, cat, valueMethod, message, fixes)

		return len(fixes) == 0
	}

	return true
}

// zeroSizedStarExpr analyzes the type in a star expression (*x) to determine if it involves a zero-sized type.
func (v *Visitor) zeroSizedStarExpr(t types.Type) (cat category, message string, valueMethod, zeroSized bool) {
	if t == nil {
		return catNone, "", false, false
	}

	if elem, vM, zS := v.base.ZeroSizedTypePointer(t); zS {
		// *... where the type of ... is a pointer to a zero-size variable.
		// p := &struct{}{}; _ = *p
		return catDeref, fmt.Sprintf("dereferencing pointer to zero-size variable of type %q", elem), vM, true
	}

	if vM, zS := v.base.ZeroSizedType(t); zS {
		// *... where ... is a zero-sized type.
		// type t struct{}; case *t:, map[]*t, []*t, A[*t]
		return catStarType, fmt.Sprintf("pointer to zero-sized type %q", t), vM, true
	}

	return catNone, "", false, false
}
