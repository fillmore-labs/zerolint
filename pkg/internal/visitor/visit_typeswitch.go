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

import "go/ast"

func (v *Visitor) visitTypeSwitch(n *ast.TypeSwitchStmt) bool {
	for _, b := range n.Body.List {
		if c, ok := b.(*ast.CaseClause); ok {
			for _, x := range c.List {
				t := v.check.TypesInfo().TypeOf(x)
				if elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(t); zeroSized {
					cM := msgFormatf(catTypeAssert, valueMethod, "type switch with pointer to zero-size variable of type %q", elem)
					fixes := v.check.RemoveStar(x)
					v.check.Report(x, cM, fixes)
				}
			}
		}
	}

	return true
}
