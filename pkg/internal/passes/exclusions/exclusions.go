// Copyright 2024-2025 Oliver Eikemeier. All Rights Reserved.
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

package exclusions

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/internal/typeutil"
)

type calc struct{ *analysis.Pass }

// CalculateExclusions performs an analysis to identify excluded types for the current pass.
//
// It works in two stages:
//
//  1. It retrieves all `excludedFact`s that were exported by the [Analyzer].
//     This covers types that were excluded directly in their defining package
//     using a `//zerolint:exclude` comment on the `type` declaration.
//     An [analysis.Fact] is used for this because it allows the exclusion information
//     to be passed between packages.
//
//  2. It scans `var` declarations in the *current* package for `//zerolint:exclude` comments
//     (e.g., `var _ another.Type`). This mechanism is necessary to exclude types from
//     external packages. We cannot export a Fact for an external type, as the analysis
//     framework requires facts to be associated with objects defined in the current package.
//     Instead, we resolve the type identifier within the current pass and add its definition
//     position directly to the result set. This local calculation must be performed by each
//     consuming analyzer.
func CalculateExclusions(ap *analysis.Pass) (set.Set[token.Pos], error) {
	c := calc{Pass: ap}

	excludedTypeDefs, err := ResultOf(c.Pass)
	if err != nil {
		return nil, err
	}

	c.addExclusions(excludedTypeDefs)

	for genDecl := range typeutil.AllDecls[*ast.GenDecl](c.Files) {
		switch genDecl.Tok { //nolint:exhaustive
		case token.TYPE:
			// Check for misplaced comments on the spec.
			c.lintSpecs(genDecl)

		case token.VAR:
			// Exclude all via "//zerolint:exclude" comment on the declaration block.
			if HasExcludeComment(genDecl.Doc) {
				// Process the spec to find the type to exclude.
				c.processExcludedValueSpec(genDecl, excludedTypeDefs)
			}

			// Check for misplaced comments on the spec.
			c.lintSpecs(genDecl)

		default: // ignore import or const
		}
	}

	// Return a [set.Set] containing the [token.Pos] of all excluded type definitions.
	return excludedTypeDefs, nil
}

// processExcludedValueSpec handles a ValueSpec within a "//zerolint:exclude" var block.
// It expects the pattern `var _ Type` and adds 'Type' to the excluded set.
func (c calc) processExcludedValueSpec(genDecl *ast.GenDecl, excludedTypeDefs set.Set[token.Pos]) {
	for _, genSpec := range genDecl.Specs {
		spec, ok := genSpec.(*ast.ValueSpec)
		if !ok { // should not happen
			log.Printf("Internal error: Expected *ast.ValueSpec, got %T (zl:xxx)", genSpec)

			continue
		}

		// Check for "_" identifier
		if len(spec.Names) != 1 || spec.Names[0].Name != "_" {
			c.Report(analysis.Diagnostic{
				Pos:     spec.Pos(),
				End:     spec.End(),
				Message: "Only one \"_\" identifier admissible (zl:com)",
			})

			continue
		}

		// Check for spec.Type and no spec.Values
		if spec.Type == nil || len(spec.Values) > 0 {
			c.Report(analysis.Diagnostic{
				Pos:     spec.Pos(),
				End:     spec.End(),
				Message: "One type name and no expressions are required (zl:com)",
			})

			continue
		}

		tv := c.TypesInfo.Types[spec.Type]

		if !tv.IsType() { // should not happen
			log.Printf("Internal error: Expected type expression, got %#v (zl:xxx)", tv)

			continue
		}

		var tn *types.TypeName

		switch t := tv.Type.(type) {
		case *types.Named:
			tn = t.Obj()

		case *types.Alias:
			tn = t.Obj()

		default:
			ts := types.TypeString(t, types.RelativeTo(c.Pkg))
			c.ReportRangef(spec, "Expected type name, got %q (zl:com)", ts)

			continue
		}

		excludedTypeDefs.Add(tn.Pos())
	}
}

// lintSpecs analyzes a GenDecl to ensure exclude comments are correctly positioned on the declaration block.
func (c calc) lintSpecs(genDecl *ast.GenDecl) {
	for _, spec := range genDecl.Specs {
		switch spec := spec.(type) {
		case *ast.TypeSpec:
			if HasExcludeComment(spec.Doc) || HasExcludeComment(spec.Comment) {
				// Comment should be on the declaration block, not the type spec.
				c.Report(analysis.Diagnostic{
					Pos:     spec.Pos(),
					End:     spec.End(),
					Message: "Exclude types with comment before the \"type\" keyword (zl:com)",
				})
			}

		case *ast.ValueSpec:
			if HasExcludeComment(spec.Doc) || HasExcludeComment(spec.Comment) {
				// Comment should be on the declaration block, not the value spec.
				c.Report(analysis.Diagnostic{
					Pos:     spec.Pos(),
					End:     spec.End(),
					Message: "Exclude types with comment before the \"var\" keyword (zl:com)",
				})
			}
		}
	}
}

// addExclusions prefills hard coded type definitions.
// For example, it ignores [runtime.Func] because pointers to this type represent opaque
// runtime-internal data, not zero-sized types the linter targets.
//
// We do it here and not as Facts:
//
// > “Some driver implementations (such as those based on Bazel and Blaze)
// do not currently apply analyzers to packages of the standard library.”
//
// https://pkg.go.dev/golang.org/x/tools/go/analysis#hdr-Modular_analysis_with_Facts
func (c calc) addExclusions(excludedTypeDefs set.Set[token.Pos]) {
	for _, pkg := range c.Pkg.Imports() {
		var typeNames []string

		switch pkg.Path() {
		case "runtime":
			typeNames = []string{"Func"}

		case "runtime/cgo":
			typeNames = []string{"Incomplete"}

		default:
			continue
		}

		scope := pkg.Scope()
		for _, name := range typeNames {
			if tn, ok := scope.Lookup(name).(*types.TypeName); ok {
				excludedTypeDefs.Add(tn.Pos())
			}
		}
	}
}
