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

package checker

import (
	"go/ast"
	"go/token"
	"go/types"
	"regexp"

	"fillmore-labs.com/zerolint/pkg/internal/set"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

// Checker has helper functions for analyzing and reporting pointers to zero-size variables.
type Checker struct {
	pass               *analysis.Pass
	Excludes, Detected set.Set[string]
	Regex              *regexp.Regexp
	ignored, seen      set.Set[token.Pos]
	cache              typeutil.Map
	Current            *ast.File
}

// typeCache stores the cached results of zero-sized and value-method checks for named types.
type typeCache struct {
	zeroSized, valueMethod bool
}

// Prepare initializes the [Checker] with the provided [analysis.Pass], preparing for new analysis.
func (c *Checker) Prepare(pass *analysis.Pass) {
	c.pass = pass

	c.ignored = set.New[token.Pos]()
	c.calcIgnored()

	if c.Excludes == nil {
		c.Excludes = set.New[string]()
	}

	c.Detected = set.New[string]()
	c.seen = set.New[token.Pos]()
}

// TypesInfo returns the type information for the current analysis pass.
func (c *Checker) TypesInfo() *types.Info {
	return c.pass.TypesInfo
}

// IgnoreType adds the specified type's declaration position to the ignored set,
// marking it as excluded from further processing.
func (c *Checker) IgnoreType(tn *types.TypeName) {
	c.ignored.Insert(tn.Pos())
}

// Print prints the syntax tree representation of the given node `x`.
func (c *Checker) Print(x any) error {
	return ast.Print(c.pass.Fset, x)
}

// calcIgnored determines ignored type definitions.
// Currently, it ignores [runtime.Func] because pointers to this type represent opaque
// runtime-internal data, not zero-sized types the linter targets.
func (c *Checker) calcIgnored() {
	lookup := map[string][]string{
		"runtime":     {"Func"},
		"runtime/cgo": {"Incomplete"},
	}

	for _, i := range c.pass.Pkg.Imports() {
		if names, ok := lookup[i.Path()]; ok {
			for _, name := range names {
				if obj := i.Scope().Lookup(name); obj != nil {
					if tn, ok := obj.(*types.TypeName); ok {
						c.IgnoreType(tn)
					}
				}
			}
		}
	}
}
