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
	"strings"

	"fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
)

// visitTypeSpec analyzes type declarations to detect if they explicitly declare types as pointers to zero-sized types.
func (v *Visitor) visitTypeSpec(n *ast.TypeSpec) bool {
	if isCtype(n) {
		return false // cgo types are often opaque
	}

	t := v.Diag.TypesInfo().TypeOf(n.Type)
	if elem, valueMethod, zeroSized := v.Check.ZeroSizedTypePointer(t); zeroSized {
		cM := msg.Formatf(msg.CatTypeDeclaration, valueMethod,
			"type declaration to pointer to zero-sized type %q", elem)
		fixes := v.removeStar(n.Type)
		v.Diag.Report(n, cM, fixes)
	}

	return true
}

// isCtype is a heuristic to avoid issues with cgo types like _Ctype_struct_foo.
func isCtype(ts *ast.TypeSpec) bool {
	const cgoTypePrefix = "_Ctype_"

	return strings.HasPrefix(ts.Name.Name, cgoTypePrefix)
}
