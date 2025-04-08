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
	"strings"
)

// visitStructType analyzes ....
func (v *Visitor) visitStructType(n *ast.StructType) bool {
	for _, field := range n.Fields.List {
		if !v.full && len(field.Names) > 0 {
			continue
		}
		t := v.Pass.TypesInfo.TypeOf(field.Type)
		if elem, ok := v.zeroSizedTypePointer(t); ok {
			var message string
			switch len(field.Names) {
			case 0:
				message = fmt.Sprintf("embedded pointer to zero-sized type %q (ZL05)", elem)
			case 1:
				message = fmt.Sprintf("field %s points to zero-sized type %q (ZL06)", field.Names[0].Name, elem)
			default:
				names := make([]string, len(field.Names))
				for i, n := range field.Names {
					names[i] = n.Name
				}
				message = fmt.Sprintf("fields %s point to zero-sized type %q (ZL06)", strings.Join(names, ", "), elem)
			}
			fixes := v.removeStar(field.Type)
			v.report(field, message, fixes)
		}
	}

	return true
}
