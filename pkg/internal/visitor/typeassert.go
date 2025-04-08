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
)

// visitTypeAssert analyzes type assertions x.(T).
func (v *Visitor) visitTypeAssert(n *ast.TypeAssertExpr) bool {
	x := n.Type
	if x == nil { // type switch
		return true
	}

	t := v.base.TypesInfo().TypeOf(x)

	elem, valueMethod, zeroSized := v.base.ZeroSizedTypePointer(t)
	if !zeroSized {
		return true
	}

	message := fmt.Sprintf("type assert to pointer to zero-size variable of type %q", elem)
	fixes := v.base.RemoveStar(x)
	v.base.Report(n, catTypeAssert, valueMethod, message, fixes)

	return true
}
