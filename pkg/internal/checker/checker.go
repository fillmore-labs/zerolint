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

	"fillmore-labs.com/zerolint/pkg/internal/filter"
	"fillmore-labs.com/zerolint/pkg/internal/passes/excluded"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

// Checker provides helper functions for analyzing and reporting pointers to zero-sized types.
type Checker struct {
	pass               *analysis.Pass
	Excludes, Detected set.Set[string]

	// Filter for zero-sized checks, used in [Checker.ZeroSizedType].
	Regex *regexp.Regexp

	// Tracks StarExpr positions that have already been processed to avoid duplicate diagnostics or fixes.
	seenStars set.Set[token.Pos]

	// Cached results of zero-sized checks, used by [Checker.ZeroSizedType] to optimize repeated lookups.
	cache typeutil.Map

	// Imports of the currently processed file, used by [Checker.Qualifier].
	CurrentImports []*ast.ImportSpec

	// Type definitions excluded via `//nolint:zerolint` directives.
	ExcludedTypeDefs filter.Filter
}

// New creates and initializes a [Checker] instance using the provided [analysis.Pass].
func New(pass *analysis.Pass) *Checker {
	c := &Checker{}
	c.Prepare(pass)

	return c
}

// Prepare initializes the [Checker] with the provided [analysis.Pass], preparing for new analysis.
func (c *Checker) Prepare(pass *analysis.Pass) {
	c.pass = pass

	if c.Excludes == nil {
		c.Excludes = set.New[string]()
	}

	c.Detected = set.New[string]()
	c.seenStars = set.New[token.Pos]()

	if excludedTypeDefs, ok := pass.ResultOf[excluded.Analyzer].(filter.Filter); ok {
		c.ExcludedTypeDefs = excludedTypeDefs
	}
}

// TypesInfo returns the type information for the current analysis pass.
func (c *Checker) TypesInfo() *types.Info {
	return c.pass.TypesInfo
}

// Print outputs the syntax tree representation of the given AST node `n` to standard output.
// This can be useful for debugging purposes.
func (c *Checker) Print(n any) error {
	return ast.Print(c.pass.Fset, n)
}
