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
	"go/token"
)

// visitUnary checks expressions in form &x.
func (v *Visitor) visitUnary(n *ast.UnaryExpr) bool {
	if n.Op != token.AND {
		return true
	}

	// &...
	t := v.check.TypesInfo().TypeOf(n.X)

	valueMethod, zeroSized := v.check.ZeroSizedType(t)
	if !zeroSized {
		return true
	}

	cM := msgFormatf(catAddress, valueMethod, "address of zero-size variable of type %q", t)
	fixes := v.check.RemoveOp(n, n.X)
	v.check.Report(n, cM, fixes)

	return len(fixes) == 0
}
