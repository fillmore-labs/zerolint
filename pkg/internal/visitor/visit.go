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
)

// visit is the central visitor function called by `inspector.Nodes`.
// It receives AST nodes during traversal and dispatches them to specific
// `visit*` methods based on the node type for detailed analysis.
func (v *Visitor) visit(n ast.Node, push bool) (proceed bool) { //nolint:cyclop
	if !push {
		return true // Only process nodes when entering
	}

	switch n := n.(type) {
	case *ast.StarExpr:
		return v.visitStar(n)

	case *ast.UnaryExpr:
		return v.visitUnary(n)

	case *ast.BinaryExpr:
		return v.visitBinary(n)

	case *ast.CallExpr:
		return v.visitCall(n)

	case *ast.File:
		return v.visitFile(n)

	case *ast.FuncDecl:
		return v.visitFuncDecl(n)

	case *ast.FuncLit:
		return v.visitFuncLit(n)

	case *ast.TypeAssertExpr:
		return v.visitTypeAssert(n)

	case *ast.StructType:
		return v.visitStructType(n)

	case *ast.FuncType:
		return v.visitFuncType(n)

	case *ast.ValueSpec:
		return v.visitValueSpec(n)

	case *ast.TypeSpec:
		return v.visitTypeSpec(n)

	default:
		return true
	}
}
