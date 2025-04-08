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
	"strings"
)

// visitTypeSpec analyzes type declarations to detect if they explicitly declare types as pointers to zero-sized types.
func (v *Visitor) visitTypeSpec(n *ast.TypeSpec) bool {
	if isCtype(n.Name.Name) {
		return false // cgo types are often opaque
	}

	t := v.check.TypesInfo().TypeOf(n.Type)
	if elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(t); zeroSized {
		cM := msgFormatf(catTypeDeclaration, valueMethod,
			"type declaration to pointer to zero-sized type %q", elem)
		fixes := v.check.RemoveStar(n.Type)
		v.check.Report(n, cM, fixes)
	}

	return true
}

func isCtype(name string) bool {
	// Heuristic to avoid issues with cgo types like _Ctype_struct_foo
	return strings.HasPrefix(name, "_Ctype_")
}
