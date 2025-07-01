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
)

// visitTypeAssert analyzes type assertions x.(T).
func (v *visitor) visitTypeAssert(n *ast.TypeAssertExpr) bool {
	x := n.Type
	if x == nil { // type switch
		return true
	}

	t := v.diag.TypesInfo().TypeOf(x)
	if elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(t); zeroSized {
		cM := msg.Formatf(msg.CatTypeAssert, valueMethod, "type assert to pointer to zero-size variable of type %q", elem)
		fixes := v.removeStar(x)
		v.diag.Report(n, cM, fixes)
	}

	return true
}
