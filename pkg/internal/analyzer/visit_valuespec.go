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

// visitValueSpec analyzes variable declarations (`var` or `const` specs)
// to detect if they explicitly declare variables as pointers to zero-sized types.
func (v *Visitor) visitValueSpec(n *ast.ValueSpec) bool {
	if n.Type == nil {
		return true
	}

	t, ok := v.Diag.TypesInfo().Types[n.Type]
	if !ok { // should not happen
		v.Diag.LogErrorf(n, "Can't find variable declaration type")

		return true
	}

	elem, valueMethod, zeroSized := v.Check.ZeroSizedTypePointer(t.Type)
	if !zeroSized {
		return true
	}

	for _, name := range n.Names {
		const cgoValuePrefix = "_cgo"

		if strings.HasPrefix(name.Name, cgoValuePrefix) {
			return false
		}
	}

	cM := msg.FormatMessage(msg.Value{}, elem, valueMethod, n.Names)
	fixes := v.removeStar(n.Type)
	v.Diag.Report(n, cM, fixes)

	return true
}
