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

// visitStar checks expressions in form *x.
func (v *Visitor) visitStar(n *ast.StarExpr) bool {
	t := v.Pass.TypesInfo.TypeOf(n.X)
	message, ok := v.zeroSizedStarExpr(t)
	if !ok {
		return true
	}

	fixes := v.removeOp(n, n.X)
	v.report(n, message, fixes)

	return fixes == nil
}

func (v *Visitor) zeroSizedStarExpr(t types.Type) (message string, ok bool) {
	if t == nil {
		return "", false
	}

	if p, ok := t.Underlying().(*types.Pointer); ok {
		if !v.zeroSizedType(p.Elem()) {
			return "", false
		}

		// *... where the type of ... is a pointer to a zero-size variable.
		// p := &struct{}{}; _ = *p
		return fmt.Sprintf("pointer to zero-size variable of type %q", p.Elem()), true
	}

	if !v.zeroSizedType(t) {
		return "", false
	}

	// *... where ... is a zero-sized type.
	// type t struct{}; var _ *t
	return fmt.Sprintf("pointer to zero-sized type %q", t), true
}
