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
	"errors"
	"go/ast"

	"fillmore-labs.com/zerolint/pkg/zerolint/level"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// ErrNoInspectorResult is returned when the ast inspector is missing.
var ErrNoInspectorResult = errors.New("zerolint: inspector result missing")

// Run performs the actual analysis on the provided [analysis.Pass].
func (v *Visitor) Run(pass *analysis.Pass) (any, error) {
	v.check.Prepare(pass)

	in, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspectorResult
	}

	nodes := nodeFilter(v.level)
	in.Nodes(nodes, v.visit)

	return any(nil), nil
}

// nodeFilter determines which AST node types to inspect based on the Visitor's configuration
// (e.g., `level` flag).
func nodeFilter(lvl level.LintLevel) []ast.Node {
	var nodes []ast.Node

	switch {
	case lvl.AtLeast(level.Full):
		// Node types are included at `full` level to perform a complete analysis.
		nodes = append(nodes,
			// keep-sorted start
			filter((*Visitor).visitFuncType),
			filter((*Visitor).visitTypeAssert),
			filter((*Visitor).visitTypeSwitch),
			filter((*Visitor).visitUnary),
			// keep-sorted end
		)

		fallthrough

	case lvl.AtLeast(level.Extended):
		// Less node types are included at `extended` level to perform a more lenient analysis.
		nodes = append(nodes,
			// keep-sorted start
			filter((*Visitor).visitFuncLit),
			filter((*Visitor).visitValueSpec),
			// keep-sorted end
		)

		fallthrough

	default:
		// Basic analysis at the default level.
		nodes = append(nodes,
			// keep-sorted start
			filter((*Visitor).visitBinary),
			filter((*Visitor).visitCall),
			filter((*Visitor).visitFile),
			filter((*Visitor).visitFuncDecl),
			filter((*Visitor).visitStar),
			filter((*Visitor).visitStructType),
			filter((*Visitor).visitTypeSpec),
			// keep-sorted end
		)
	}

	return nodes
}

// filter is a helper function used in nodeFilter to obtain a zero-value
// instance of a specific ast.Node implementation.
// The function argument is not actually called, but used by Go's type inference to determine the type `N`.
// For example, for `(*Visitor).visitCall` which has the signature
// type `func(*Visitor, *ast.CallExpr) bool`, `N` is inferred as `*ast.CallExpr`.
// This allows nodeFilter to specify node types by referencing their corresponding visit methods.
func filter[N ast.Node](func(*Visitor, N) bool) ast.Node {
	var n N

	return n
}

// visit is the central visitor function called by `inspector.Nodes`.
// It receives AST nodes during traversal and dispatches them to specific
// `visit*` methods based on the node type for detailed analysis.
func (v *Visitor) visit(n ast.Node, push bool) (proceed bool) { //nolint:cyclop
	if !push {
		return true // Only process nodes when entering
	}

	switch n := n.(type) {
	// keep-sorted start
	case *ast.BinaryExpr:
		return v.visitBinary(n)
	case *ast.CallExpr:
		return v.visitCall(n)
	case *ast.File:
		return v.visitFile(n)
	case *ast.FuncDecl:
		return v.visitFuncDecl(n)
	case *ast.StarExpr:
		return v.visitStar(n)
	case *ast.StructType:
		return v.visitStructType(n)
	case *ast.TypeSpec:
		return v.visitTypeSpec(n)
	// keep-sorted end

	// keep-sorted start
	case *ast.FuncLit:
		return v.visitFuncLit(n)
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

	default:
		return true
	}
}
