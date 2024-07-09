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

func (v Visitor) visitStar(x *ast.StarExpr) bool {
	t := v.TypesInfo.Types[x.X].Type
	var message string

	p, ok := t.Underlying().(*types.Pointer)
	switch {
	case ok:
		e := p.Elem()
		if !v.isZeroSizeType(p.Elem()) {
			return true
		}
		message = fmt.Sprintf("pointer to zero-size variable of type %q", e.String())

	case v.isZeroSizeType(t):
		message = fmt.Sprintf("pointer to zero-size type %q", t.String())

	default:
		return true
	}

	fixes := removeOp(x, x.X)
	v.report(x, message, fixes)

	return fixes == nil
}
