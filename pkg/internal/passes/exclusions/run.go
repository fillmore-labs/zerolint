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
	"strings"

	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/pkg/internal/typeutil"
)

type pass struct{ *analysis.Pass }

// run performs the analysis to identify excluded types that are defined in the current package.
// It exports an [excludedFact] for each newly identified excluded type.
//
// Note that this pass only handles exclusions on `type` declarations (e.g., `//zerolint:exclude type T ...`).
// Exclusions on `var` declarations (e.g., `//zerolint:exclude var _ another.Type`) must be handled
// by the consuming analyzer via the [CalculateExclusions] function, as facts cannot be exported
// for objects defined in other packages.
func run(ap *analysis.Pass) (any, error) {
	p := pass{Pass: ap}

	for decl := range typeutil.AllDecls[*ast.GenDecl](p.Files) {
		if decl.Tok == token.TYPE {
			p.processTypeDecl(decl)
		}
	}

	return p.newResult(), nil
}

func (p pass) processTypeDecl(genDecl *ast.GenDecl) {
	// Exclude all via "//zerolint:exclude" comment on the declaration block.
	excludeAll := HasExcludeComment(genDecl.Doc)

	for _, genSpec := range genDecl.Specs {
		spec, ok := genSpec.(*ast.TypeSpec)
		if !ok { // should not happen
			log.Printf("Internal error: Expected *ast.TypeSpec, got %T (zl:xxx)", genSpec)

			continue
		}

		// A type is excluded if it's in a `//zerolint:exclude` block
		// or if it's a CGo-generated type.
		if excludeAll || isCtype(spec) {
			def := p.TypesInfo.Defs[spec.Name]

			if tn, ok := def.(*types.TypeName); ok {
				p.excludeType(tn)
			} else { // should not happen
				log.Printf("Internal error: Expected *types.TypeName, got %T (zl:xxx)", def)
			}
		}
	}
}

// isCtype is a heuristic to avoid issues with cgo types like _Ctype_struct_foo.
func isCtype(ts *ast.TypeSpec) bool {
	const cgoTypePrefix = "_Ctype_"

	return strings.HasPrefix(ts.Name.Name, cgoTypePrefix)
}
