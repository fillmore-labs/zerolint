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
	"bytes"
	"go/ast"
	"go/format"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

func removeOp(n ast.Node, x ast.Expr) []analysis.SuggestedFix {
	var fixes []analysis.SuggestedFix
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), x); err == nil {
		edit := analysis.TextEdit{
			Pos:     n.Pos(),
			End:     n.End(),
			NewText: buf.Bytes(),
		}
		fixes = []analysis.SuggestedFix{
			{Message: "change to pure type", TextEdits: []analysis.TextEdit{edit}},
		}
	}

	return fixes
}

func (v Visitor) report(pos, end token.Pos, message string, fixes []analysis.SuggestedFix) {
	v.Report(analysis.Diagnostic{
		Pos:            pos,
		End:            end,
		Category:       "zero-size",
		Message:        message,
		URL:            "https://pkg.go.dev/fillmore-labs.com/zerolint",
		SuggestedFixes: fixes,
	})
}
