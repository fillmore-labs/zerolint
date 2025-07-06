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

// checkFieldList examines a field list (struct fields, function parameters, or results) for pointers
// to zero-sized types.
//
// If skipNamed is true, it will skip fields with names (useful for examining only embedded types in structs).
// The msgFormatter provides context-specific messages based on whether the current context involves struct fields,
// function parameters, or function results.
func (v *Visitor) checkFieldList(n *ast.FieldList, skipNamed bool, formatter msg.Formatter) {
	for _, field := range n.List {
		if skipNamed && len(field.Names) > 0 {
			continue // check only embedded types
		}

		t := v.Diag.TypesInfo().TypeOf(field.Type)
		if elem, valueMethod, zeroSized := v.Check.ZeroSizedTypePointer(t); zeroSized {
			cM := msg.FormatMessage(formatter, elem, valueMethod, field.Names)
			fixes := v.removeStar(field.Type)
			v.Diag.Report(field, cM, fixes)
		}
	}
}
