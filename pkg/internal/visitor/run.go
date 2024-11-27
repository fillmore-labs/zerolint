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
	"slices"

	"fillmore-labs.com/zerolint/pkg/internal/set"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Run runs the analysis.
func Run(logger *log.Logger, visitor Visitor, zeroTrace, basic, generated bool) {
	in, ok := visitor.Pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		logger.Fatal("inspector result missing")
	}

	visitor.Excludes.Insert("runtime.Func")
	visitor.Detected = set.New[string]()

	types, f := visitFunc(visitor, basic, generated)
	in.Nodes(types, f)

	if zeroTrace && len(visitor.Detected) > 0 {
		logger.Printf("found zero-sized types in %q:\n", visitor.Pass.Pkg.Path())
		names := visitor.Detected.Elements()
		slices.Sort(names)
		for _, name := range names {
			logger.Printf("- %s\n", name)
		}
	}
}

// visitFunc determines parameters and function to call for inspector.Nodes.
func visitFunc(visitor Visitor, basic, generated bool) ([]ast.Node, func(ast.Node, bool) bool) {
	types := make([]ast.Node, 0, 5) //nolint:mnd
	var f func(ast.Node, bool) bool

	types = append(types, (*ast.BinaryExpr)(nil), (*ast.CallExpr)(nil))
	if !generated {
		types = append(types, (*ast.File)(nil))
	}

	if basic {
		types = append(types, (*ast.FuncDecl)(nil))
		f = visitor.visitBasic
	} else {
		types = append(types, (*ast.StarExpr)(nil), (*ast.UnaryExpr)(nil))
		f = visitor.visit
	}

	return types, f
}
