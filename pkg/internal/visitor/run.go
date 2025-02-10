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
	"go/token"
	"iter"
	"log"
	"strings"

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
	Full      bool
	Generated bool
}

// Visitor is an AST visitor for analyzing usage of pointers to zero-size variables.
type Visitor struct {
	Pass            *analysis.Pass
	excludes        set.Set[string]
	detected        set.Set[string]
	gen             set.Set[*ast.File]
	ignored         set.Set[token.Pos]
	full, generated bool
}

// New returns a [Visitor] configured with [Options].
func New(opt Options) *Visitor {
	return &Visitor{
		excludes:  opt.Excludes,
		full:      opt.Full,
		generated: opt.Generated,
	}
}

// HasDetected tells whether zero-sized types have been detected.
func (v *Visitor) HasDetected() bool {
	return len(v.detected) > 0
}

// AllDetected returns an iterator over the detected zero-sized types in alphabetical order.
func (v *Visitor) AllDetected() iter.Seq[string] {
	return set.AllSorted(v.detected)
}

// ErrNoInspectorResult is returned when the ast inspetor is missing.
var ErrNoInspectorResult = errors.New("zerolint: inspector result missing")

const zerolintExclude = "zerolint:exclude"

// Run runs the analysis.
func (v *Visitor) Run(pass *analysis.Pass) (any, error) {
	v.Pass = pass
	in, ok := v.Pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspectorResult
	}

	v.calcIgnored(in)

	if v.excludes == nil {
		v.excludes = set.New[string]()
	}
	v.excludes.Insert("runtime.Func")
	v.detected = set.New[string]()

	types, f := v.selectVisitFunc()
	in.Nodes(types, f)

	return any(nil), nil
}

func (v *Visitor) calcIgnored(in *inspector.Inspector) {
	v.gen = set.New[*ast.File]()
	v.ignored = set.New[token.Pos]()
	var isGen bool
	for n := range in.PreorderSeq((*ast.File)(nil), (*ast.TypeSpec)(nil)) {
		switch n := n.(type) {
		case *ast.File:
			isGen = ast.IsGenerated(n)
			if isGen {
				v.gen.Insert(n)
			}

		case *ast.TypeSpec:
			if isGen {
				v.ignored.Insert(n.Name.NamePos)

				continue
			}

			if group := n.Comment; group != nil {
				for _, c := range group.List {
					if strings.Contains(c.Text, zerolintExclude) {
						v.ignored.Insert(n.Name.NamePos)

						break
					}
				}
			}
		}
	}
}

// visitFunc determines parameters and function to call for inspector.Nodes.
func (v *Visitor) selectVisitFunc() ([]ast.Node, func(n ast.Node, push bool) (proceed bool)) {
	types := make([]ast.Node, 0, 5) //nolint:mnd

	types = append(types, (*ast.BinaryExpr)(nil), (*ast.CallExpr)(nil))
	if !v.generated {
		types = append(types, (*ast.File)(nil))
	}

	if v.full {
		types = append(types, (*ast.StarExpr)(nil), (*ast.UnaryExpr)(nil))
	}

	return types, v.visit
}
