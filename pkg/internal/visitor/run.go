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
	"iter"
	"log"

	"fillmore-labs.com/zerolint/pkg/analyzer/level"
	"fillmore-labs.com/zerolint/pkg/internal/base"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Options defines configurable parameters for the linter.
type Options struct {
	Logger    *log.Logger
	Excludes  set.Set[string]
	ZeroTrace bool
	Level     level.LintLevel
	Generated bool
}

// Visitor is an AST visitor for analyzing the usage of pointers to zero-size variables.
// It identifies various patterns where such pointers might be used unnecessarily.
type Visitor struct {
	base      base.Base
	level     int
	generated bool
}

// New creates a new [Visitor] configured with the provided [Options].
func New(opt Options) *Visitor {
	return &Visitor{
		base: base.Base{
			Excludes: opt.Excludes,
		},
		level:     int(opt.Level),
		generated: opt.Generated,
	}
}

// HasDetected tells whether any zero-sized types have been detected during analysis.
func (v *Visitor) HasDetected() bool {
	return len(v.base.Detected) > 0
}

// AllDetected returns a sorted iterator over all detected zero-sized types.
func (v *Visitor) AllDetected() iter.Seq[string] {
	return set.AllSorted(v.base.Detected)
}

// ErrNoInspectorResult is returned when the ast inspetor is missing.
var ErrNoInspectorResult = errors.New("zerolint: inspector result missing")

// Run performs the actual analysis on the provided [analysis.Pass].
func (v *Visitor) Run(pass *analysis.Pass) (any, error) {
	v.base.Prepare(pass)

	if v.base.Excludes == nil {
		v.base.Excludes = set.New[string]()
	}

	nodes, f := v.selectVisitFunc()

	in, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspectorResult
	}

	in.Nodes(nodes, f)

	return any(nil), nil
}

// selectVisitFunc determines which AST node types to inspect based on the Visitor's configuration
// (e.g., `level`, `generated` flags) and returns the appropriate visitor function for the inspector.
func (v *Visitor) selectVisitFunc() ([]ast.Node, func(n ast.Node, push bool) (proceed bool)) {
	const maxNodes = 12
	nodes := make([]ast.Node, 0, maxNodes)

	nodes = append(nodes,
		(*ast.BinaryExpr)(nil),
		(*ast.CallExpr)(nil),
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.StructType)(nil),
		(*ast.TypeSpec)(nil),
	)

	if v.level > 1 {
		// More node nodes are included when `level` is true to perform a more thorough analysis.
		nodes = append(nodes,
			(*ast.FuncLit)(nil),
			(*ast.FuncType)(nil),
			(*ast.StarExpr)(nil),
			(*ast.TypeAssertExpr)(nil),
			(*ast.UnaryExpr)(nil),
			(*ast.ValueSpec)(nil),
		)
	}

	return nodes, v.visit
}
