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
	"log"

	"fillmore-labs.com/zerolint/pkg/set"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Run struct {
	Visitor
	ZeroTrace bool
	Basic     bool
}

// Visitor is an AST visitor for analyzing usage of pointers to zero-size variables.
type Visitor struct {
	*analysis.Pass
	Excludes set.Set[string]
	Detected set.Set[string]
}

// Run runs the analysis.
func (v Run) Run() {
	v.Excludes.Insert("runtime.Func")
	v.Detected = set.New[string]()

	if in, ok := v.ResultOf[inspect.Analyzer].(*inspector.Inspector); ok {
		in.Nodes(v.visitFunc())
	} else {
		log.Fatal("inspector result missing")
	}

	if v.ZeroTrace {
		for name := range v.Detected {
			log.Printf("found zero-sized type %q", name)
		}
	}
}

// visitFunc determines parameters and function to call for inspector.Nodes.
func (v Run) visitFunc() ([]ast.Node, func(ast.Node, bool) bool) {
	if v.Basic {
		return []ast.Node{
				(*ast.BinaryExpr)(nil),
				(*ast.CallExpr)(nil),
				(*ast.FuncDecl)(nil),
			},
			v.visitBasic
	}

	return []ast.Node{
			(*ast.StarExpr)(nil),
			(*ast.UnaryExpr)(nil),
			(*ast.BinaryExpr)(nil),
			(*ast.CallExpr)(nil),
		},
		v.visit
}

// visitBasic is the main functions called by inspector.Nodes for basic analysis.
func (v Visitor) visitBasic(n ast.Node, push bool) bool {
	if !push {
		return true
	}

	switch x := n.(type) {
	case *ast.BinaryExpr:
		return v.visitBinary(x)

	case *ast.CallExpr:
		return v.visitCallBasic(x)

	case *ast.FuncDecl:
		return v.visitFunc(x)

	default:
		return true
	}
}

// visit is the main functions called by inspector.Nodes for full analysis.
func (v Visitor) visit(n ast.Node, push bool) bool {
	if !push {
		return true
	}

	switch x := n.(type) {
	case *ast.StarExpr:
		return v.visitStar(x)

	case *ast.UnaryExpr:
		return v.visitUnary(x)

	case *ast.BinaryExpr:
		return v.visitBinary(x)

	case *ast.CallExpr:
		return v.visitCall(x)

	default:
		return true
	}
}
