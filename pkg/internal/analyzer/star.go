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

	"golang.org/x/tools/go/analysis"
)

// removeStar suggests a fix that removes the star ('*') operator from an expression, if possible.
func (v *visitor) removeStar(x ast.Expr) []analysis.SuggestedFix {
	if n, ok := ast.Unparen(x).(*ast.StarExpr); ok {
		v.ignoreStar(n)

		return v.diag.RemoveOp(n, n.X)
	}
	// We could handle aliases with `*ast.Ident` here, but assume we already fixed the alias definition.

	return nil
}

// ignoreStar ignores the star expression in further processing.
func (v *visitor) ignoreStar(n *ast.StarExpr) {
	v.seenStars.Insert(n.Pos())
}

// starSeen checks if a given star expression has already been processed by verifying its position in the seen set.
func (v *visitor) starSeen(n *ast.StarExpr) bool {
	return v.seenStars.Has(n.Pos())
}
