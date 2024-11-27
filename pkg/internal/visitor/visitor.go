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

	"fillmore-labs.com/zerolint/pkg/internal/set"
	"golang.org/x/tools/go/analysis"
)

// Visitor is an AST visitor for analyzing usage of pointers to zero-size variables.
type Visitor struct {
	Pass     *analysis.Pass
	Excludes set.Set[string]
	Detected set.Set[string]
}

// visitBasic is the main functions called by inspector.Nodes for basic analysis.
func (v Visitor) visitBasic(x ast.Node, push bool) bool {
	if !push {
		return true
	}

	switch x := x.(type) {
	case *ast.BinaryExpr:
		return v.visitBinary(x)

	case *ast.CallExpr:
		return v.visitCallBasic(x)

	case *ast.FuncDecl:
		return v.visitFunc(x)

	case *ast.File:
		return v.visitFile(x)

	default:
		return true
	}
}

// visit is the main functions called by inspector.Nodes for full analysis.
func (v Visitor) visit(x ast.Node, push bool) bool {
	if !push {
		return true
	}

	switch x := x.(type) {
	case *ast.StarExpr:
		return v.visitStar(x)

	case *ast.UnaryExpr:
		return v.visitUnary(x)

	case *ast.BinaryExpr:
		return v.visitBinary(x)

	case *ast.CallExpr:
		return v.visitCall(x)

	case *ast.File:
		return v.visitFile(x)

	default:
		return true
	}
}
