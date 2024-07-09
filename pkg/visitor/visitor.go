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

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Visitor struct {
	*analysis.Pass
	Excludes  map[string]struct{}
	Detected  map[string]struct{}
	ZeroTrace bool
	Basic     bool
}

func (v Visitor) Run() {
	in, ok := v.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		log.Fatal("inspector result missing")
	}

	v.Excludes["runtime.Func"] = struct{}{}
	v.Detected = make(map[string]struct{})

	if v.Basic {
		types := []ast.Node{
			(*ast.BinaryExpr)(nil),
			(*ast.CallExpr)(nil),
		}
		in.Nodes(types, v.visitBasic)
	} else {
		types := []ast.Node{
			(*ast.StarExpr)(nil),
			(*ast.UnaryExpr)(nil),
			(*ast.BinaryExpr)(nil),
			(*ast.CallExpr)(nil),
		}
		in.Nodes(types, v.visit)
	}

	if v.ZeroTrace {
		for name := range v.Detected {
			log.Printf("found zero-size type %q", name)
		}
	}
}

func (v Visitor) visitBasic(n ast.Node, push bool) bool {
	if !push {
		return true
	}

	switch x := n.(type) {
	case *ast.BinaryExpr:
		return v.visitBinary(x)

	case *ast.CallExpr:
		return v.visitCallBasic(x)

	default:
		return true
	}
}

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
