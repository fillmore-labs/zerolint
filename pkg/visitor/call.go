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

func (v Visitor) visitCall(x *ast.CallExpr) bool {
	if !isSingleNil(x.Args) || !v.isZeroPointer(x.Fun) {
		return true
	}

	var fixes []analysis.SuggestedFix
	s, ok := unwrap(x.Fun).(*ast.StarExpr)
	if ok {
		var buf bytes.Buffer
		if err := format.Node(&buf, token.NewFileSet(), s.X); err == nil {
			buf.WriteString("{}")
			edit := analysis.TextEdit{
				Pos:     x.Pos(),
				End:     x.End(),
				NewText: buf.Bytes(),
			}
			fixes = []analysis.SuggestedFix{
				{Message: "change to pure type", TextEdits: []analysis.TextEdit{edit}},
			}
		}
	}

	const message = "cast to pointer to zero-size variable"
	v.report(x.Pos(), x.End(), message, fixes)

	return fixes == nil
}

func isSingleNil(args []ast.Expr) bool {
	if len(args) != 1 {
		return false
	}
	i, ok := args[0].(*ast.Ident)

	return ok && i.Name == "nil"
}

func unwrap(e ast.Expr) ast.Expr {
	x := e
	for {
		p, ok := x.(*ast.ParenExpr)
		if !ok {
			return x
		}
		x = p.X
	}
}
