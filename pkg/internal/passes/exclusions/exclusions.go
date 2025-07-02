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
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/pkg/internal/set"
)

type calc analysis.Pass

// CalculateExclusions performs the analysis to identify excluded types.
// It retrieves exclusion facts from the dedicated exclusions pass and
// adds types excluded via the "//zerolint:exclude" comment on "var _ Type" declarations.
func CalculateExclusions(pass *analysis.Pass) (set.Set[token.Pos], error) {
	c := (*calc)(pass)

	excludedTypeDefs, err := ResultOf(pass)
	if err != nil {
		return nil, err
	}

	for genDecl := range AllDecl[*ast.GenDecl](c.Files) {
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
func (c *calc) processExcludedValueSpec(genDecl *ast.GenDecl, excludedTypeDefs set.Set[token.Pos]) {
	for _, genSpec := range genDecl.Specs {
		spec, ok := genSpec.(*ast.ValueSpec)
		if !ok { // should not happen
			log.Printf("Internal error: Expected *ast.ValueSpec, got %T (zl:xxx)", genSpec)

			continue
		}

		id := spec.Names[0]

		// Check for "_" identifier
		if len(spec.Names) != 1 || id.Name != "_" {
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

		// Look up the type definition
		def, ok := c.TypesInfo.Defs[id]
		if !ok { // should not happen
			log.Printf("Internal error: can't find value type definition (zl:xxx)")

			continue
		}

		var tn *types.TypeName
		switch t := def.Type().(type) {
		case *types.Named:
			tn = t.Obj()

		case *types.Alias:
			tn = t.Obj()

		default:
			ts := types.TypeString(t, types.RelativeTo(c.Pkg))
			c.Report(analysis.Diagnostic{
				Pos:     spec.Pos(),
				End:     spec.End(),
				Message: fmt.Sprintf("Expected type name, got %q (zl:com)", ts),
			})

			continue
		}

		excludedTypeDefs.Insert(tn.Pos())
	}
}

// lintSpecs analyzes a GenDecl to ensure exclude comments are correctly positioned on the declaration block.
func (c *calc) lintSpecs(genDecl *ast.GenDecl) {
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
