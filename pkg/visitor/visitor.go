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
	Generated bool
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

	if v.ZeroTrace && len(v.Detected) > 0 {
		log.Printf("found zero-sized types in %q:\n", v.Pass.Pkg.Path())
		for name := range v.Detected {
			log.Printf("- %s\n", name)
		}
	}
}

// visitFunc determines parameters and function to call for inspector.Nodes.
func (v Run) visitFunc() ([]ast.Node, func(ast.Node, bool) bool) {
	types := make([]ast.Node, 0, 5) //nolint:mnd
	var f func(ast.Node, bool) bool

	types = append(types, (*ast.BinaryExpr)(nil), (*ast.CallExpr)(nil))
	if !v.Generated {
		types = append(types, (*ast.File)(nil))
	}

	if v.Basic {
		types = append(types, (*ast.FuncDecl)(nil))
		f = v.visitBasic
	} else {
		types = append(types, (*ast.StarExpr)(nil), (*ast.UnaryExpr)(nil))
		f = v.visit
	}

	return types, f
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
		return !ast.IsGenerated(x)

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
		return !ast.IsGenerated(x)

	default:
		return true
	}
}
