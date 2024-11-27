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

// visitUnary checks expressions in form *x.
func (v Visitor) visitStar(x *ast.StarExpr) bool {
	t := v.Pass.TypesInfo.TypeOf(x.X)
	var message string
	switch p, ok := t.Underlying().(*types.Pointer); {
	case ok && v.zeroSizedType(p.Elem()):
		// *... where the type of ... is a pointer to a zero-size variable.
		// p := &struct{}{}; _ = *p
		message = fmt.Sprintf("pointer to zero-size variable of type %q", p.Elem())

	case !ok && v.zeroSizedType(t):
		// *... where ... is a zero-sized type.
		// type t struct{}; var _ *t
		message = fmt.Sprintf("pointer to zero-sized type %q", t)

	default:
		return true
	}

	fixes := v.removeOp(x, x.X)
	v.report(x, message, fixes)

	return fixes == nil
}
