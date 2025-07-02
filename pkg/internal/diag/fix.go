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

package diag

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// ReplaceWithZeroValue generates a suggested fix to replace a pointer expression with its zero-value representation.
func (d *Diag) ReplaceWithZeroValue(n ast.Node, t types.Type) []analysis.SuggestedFix {
	switch t.(type) {
	case *types.Named, *types.Alias, *types.Struct, *types.Array:
		// supported types

	default:
		// types with non-zero sizes
		return nil
	}

	q := Qualifier{
		Pkg: d.pass.Pkg,
	}

	if d.CurrentFile != nil {
		q.Imports = d.CurrentFile.Imports
	}

	var buf bytes.Buffer
	types.WriteType(&buf, t, q.Qualifier)

	if q.NeedsImport {
		return nil
	}

	buf.WriteString("{}")

	return suggestedFix(n, buf.Bytes(), "replace by zero value")
}

// RemoveOp suggests a fix that removes an unary operator ('*' or '&') from an expression.
func (d *Diag) RemoveOp(n ast.Node, x ast.Expr) []analysis.SuggestedFix {
	var buf bytes.Buffer
	if err := format.Node(&buf, d.pass.Fset, x); err != nil {
		// should not happen
		d.LogErrorf(n, "Unexpected error during remove formatting: %v", err)

		return nil
	}

	return suggestedFix(n, buf.Bytes(), "remove operator")
}

// MakePure suggests a fix that replaces expressions allocating or casting to a pointer
// of a zero-sized type, such as (*T)(nil) or new(T), with a value literal T{}.
// This promotes using zero-sized types directly by value rather than through pointers.
func (d *Diag) MakePure(n ast.Node, x ast.Expr) []analysis.SuggestedFix {
	var buf bytes.Buffer
	if err := format.Node(&buf, d.pass.Fset, x); err != nil {
		// should not happen
		d.LogErrorf(n, "Unexpected error during pure formatting: %v", err)

		return nil
	}

	buf.WriteString("{}")

	return suggestedFix(n, buf.Bytes(), "change to pure type")
}

// suggestedFix returns a slice of SuggestedFix containing a single fix with the specified message and text edit.
// The text edit replaces the content of the given ast.Node with the provided newText.
func suggestedFix(n ast.Node, newText []byte, message string) []analysis.SuggestedFix {
	edit := analysis.TextEdit{
		Pos:     n.Pos(),
		End:     n.End(),
		NewText: newText,
	}

	return []analysis.SuggestedFix{
		{
			Message:   message,
			TextEdits: []analysis.TextEdit{edit},
		},
	}
}
