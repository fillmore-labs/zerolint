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

	"golang.org/x/tools/go/ast/inspector"

	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// dispatch is the central visitor function called by `inspector.Nodes`.
// It receives AST nodes during traversal and dispatches them to specific
// `visit*` methods based on the node type for detailed analysis.
func (v *Visitor) dispatch(c inspector.Cursor) (proceed bool) {
	switch n := c.Node().(type) {
	// keep-sorted start
	case *ast.BinaryExpr:
		return v.visitBinary(n)
	case *ast.CallExpr:
		return v.visitCall(n)
	case *ast.File:
		return v.visitFile(n)
	// keep-sorted end

	// keep-sorted start
	case *ast.FuncDecl:
		return v.visitFuncDecl(c, n)
	case *ast.StarExpr:
		return v.visitStar(n)
	case *ast.StructType:
		return v.visitStructType(n)
	case *ast.TypeSpec:
		return v.visitTypeSpec(n)
	// keep-sorted end

	// keep-sorted start
	case *ast.FuncLit:
		return v.visitFuncLit(c, n)
	case *ast.ValueSpec:
		return v.visitValueSpec(n)
	// keep-sorted end

	// keep-sorted start
	case *ast.FuncType:
		return v.visitFuncType(n)
	case *ast.TypeAssertExpr:
		return v.visitTypeAssert(n)
	case *ast.TypeSwitchStmt:
		return v.visitTypeSwitch(n)
	case *ast.UnaryExpr:
		return v.visitUnary(n)
	// keep-sorted end

	default: // should not happen
		v.Diag.LogErrorf(n, "Unexpected dispatch type %T", n)

		return true
	}
}

// nodeFilter determines which AST node types to inspect based on the Visitor's configuration
// (e.g., `level` flag).
func (v *Visitor) nodeFilter() []ast.Node {
	var nodes []ast.Node

	switch {
	case v.Level.AtLeast(level.Full):
		// Node types are included at `full` level to perform a complete analysis.
		nodes = append(nodes,
			// keep-sorted start ignore_prefixes=nodeN,nodeC
			nodeN((*Visitor).visitFuncType),
			nodeN((*Visitor).visitTypeAssert),
			nodeN((*Visitor).visitTypeSwitch),
			nodeN((*Visitor).visitUnary),
			// keep-sorted end
		)

		fallthrough

	case v.Level.AtLeast(level.Extended):
		// Less node types are included at `extended` level to perform a more lenient analysis.
		nodes = append(nodes,
			// keep-sorted start ignore_prefixes=nodeN,nodeC
			nodeC((*Visitor).visitFuncLit),
			nodeN((*Visitor).visitValueSpec),
			// keep-sorted end
		)

		fallthrough

	case v.Level.AtLeast(level.Basic):
		// Basic analysis at the default level.
		nodes = append(nodes,
			// keep-sorted start ignore_prefixes=nodeN,nodeC
			nodeC((*Visitor).visitFuncDecl),
			nodeN((*Visitor).visitStar),
			nodeN((*Visitor).visitStructType),
			nodeN((*Visitor).visitTypeSpec),
			// keep-sorted end
		)

		fallthrough

	default:
		// Minimal analysis, not configurable.
		nodes = append(nodes,
			// keep-sorted start ignore_prefixes=nodeN,nodeC
			nodeN((*Visitor).visitBinary),
			nodeN((*Visitor).visitCall),
			nodeN((*Visitor).visitFile),
			// keep-sorted end
		)
	}

	return nodes
}

// nodeN is a helper function used in nodeFilter to obtain a zero-value
// instance of a specific ast.Node implementation.
// The function argument is not actually called, but used by Go's type inference to determine the type `N`.
// For example, for `(*Visitor).visitCall` which has the signature
// type `func(*Visitor, *ast.CallExpr) bool`, `N` is inferred as `*ast.CallExpr`.
// This allows nodeFilter to specify node types by referencing their corresponding visit methods.
func nodeN[N ast.Node](func(*Visitor, N) bool) ast.Node {
	var n N

	return n
}

// nodeC is a helper function, see [nodeN].
func nodeC[N ast.Node](func(*Visitor, inspector.Cursor, N) bool) ast.Node {
	var n N

	return n
}
